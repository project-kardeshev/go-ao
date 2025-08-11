# go-ao

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/kardeshev/go-ao)](https://goreportcard.com/report/github.com/kardeshev/go-ao)

A comprehensive Go package for interacting with the AO network, providing easy-to-use APIs for process management, message handling, and data item signing.

## Installation

```bash
go get github.com/project-kardeshev/go-ao
```

## Quick Start

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/project-kardeshev/go-ao/ao"
    "github.com/project-kardeshev/go-ao/signers"
    goarTypes "github.com/everFinance/goar/types"
)

func main() {
    // Load your Arweave wallet
    wallet, err := os.ReadFile("wallet.json")
    if err != nil {
        panic(err)
    }
    
    // Create a signer
    signer, err := signers.NewArweaveSigner(wallet)
    if err != nil {
        panic(err)
    }
    
    // Create a process client
    client := ao.NewProcessClient(
        nil, // ProcessId (nil if you plan to spawn new process)
        ao.DefaultCuUrl,
        ao.DefaultMuUrl,
        signer,
    )
    
    // Spawn a new process
    processId, result, err := client.Spawn(ao.SpawnInput{
        Module:    ao.AOSModule,
        Authority: ao.DefaultAuthority,
        Scheduler: ao.DefaultScheduler,
        Tags: []goarTypes.Tag{
            {Name: "Your-Tag", Value: "Your-Value"},
        },
    })
    
    fmt.Printf("Process spawned: %s\n", processId)
}
```

## API Reference

### ProcessClient

The `ProcessClient` is the main interface for interacting with AO processes.

#### Constructor

```go
func NewProcessClient(processId *string, cuUrl string, muUrl string, signer signers.DataItemSigner) AOClient
```

**Parameters:**
- `processId`: Optional process ID. Set to `nil` to spawn a new process, or provide an existing process ID
- `cuUrl`: Compute Unit URL (use `ao.DefaultCuUrl` for testnet)
- `muUrl`: Message Unit URL (use `ao.DefaultMuUrl` for testnet)  
- `signer`: Data item signer (Arweave or Ethereum)

#### Methods

### Spawn

Create a new AO process.

```go
func (pc *ProcessClient) Spawn(input SpawnInput) (id string, result *Result, err error)
```

**SpawnInput:**
```go
type SpawnInput struct {
    Module    string              // Process module (use ao.AOSModule for AOS)
    Authority string              // Authority address (use ao.DefaultAuthority)
    Scheduler string              // Scheduler address (use ao.DefaultScheduler)
    Tags      []goarTypes.Tag     // Additional tags
    Data      any                 // Optional data (string, number, or bytes)
    Target    string              // Optional target
    Anchor    string              // Optional anchor
}
```

**Example:**
```go
processId, result, err := client.Spawn(ao.SpawnInput{
    Module:    ao.AOSModule,
    Authority: ao.DefaultAuthority,
    Scheduler: ao.DefaultScheduler,
    Tags: []goarTypes.Tag{
        {Name: "App-Name", Value: "MyApp"},
        {Name: "App-Version", Value: "1.0.0"},
    },
})
```

### Write

Send a message to an AO process.

```go
func (pc *ProcessClient) Write(input WriteInput) (id string, result *Result, err error)
```

**WriteInput:**
```go
type WriteInput struct {
    Process string              // Ignored - always uses ProcessClient's ProcessId
    Anchor  *string            // Optional anchor (auto-generated if nil)
    Data    string             // Message data (Lua code for AOS processes)
    Tags    []goarTypes.Tag    // Message tags
}
```

**Example:**
```go
messageId, result, err := client.Write(ao.WriteInput{
    Data: `
        local message = "Hello from Go!"
        print(message)
        return message
    `,
    Tags: []goarTypes.Tag{
        {Name: "Action", Value: "Eval"},
    },
})
```

### Read

Read the current state of an AO process (dry run).

```go
func (pc *ProcessClient) Read(input DryRunInput) (*Result, error)
```

**DryRunInput:**
```go
type DryRunInput struct {
    Id     string              // Message ID
    Owner  string              // Owner address
    From   string              // Sender address  
    Anchor string              // Anchor
    Data   string              // Query data
    Tags   []goarTypes.Tag     // Query tags
}
```

**Example:**
```go
result, err := client.Read(ao.DryRunInput{
    Id:     "query-id",
    Owner:  signer.GetAddress(),
    From:   signer.GetAddress(),
    Data:   "return 'Current state: ' .. tostring(State)",
    Tags: []goarTypes.Tag{
        {Name: "Action", Value: "Read"},
    },
})
```

### Result Structure

All operations return a `Result` struct:

```go
type Result struct {
    Messages    []Message          // Output messages
    Assignments []any             // Process assignments
    Spawns      []any             // Spawned processes
    Output      map[string]any    // Process output data
    Error       any               // Process errors
    GasUsed     any               // Gas consumption
}
```

## Signers

The package supports multiple signing methods for different wallet types.

### Arweave Signer

For Arweave wallets (JWK format):

```go
import "github.com/project-kardeshev/go-ao/signers"

// Load wallet from file
wallet, err := os.ReadFile("arweave-wallet.json")
if err != nil {
    panic(err)
}

// Create signer
signer, err := signers.NewArweaveSigner(wallet)
if err != nil {
    panic(err)
}
```

### Ethereum Signer

For Ethereum private keys:

```go
import "github.com/project-kardeshev/go-ao/signers"

// Private key as hex string (with or without 0x prefix)
privateKey := "your-ethereum-private-key"

// Create signer
signer, err := signers.NewEthereumSigner(privateKey)
if err != nil {
    panic(err)
}
```

### DataItemSigner Interface

Both signers implement the `DataItemSigner` interface:

```go
type DataItemSigner interface {
    CreateAndSignDataItem(data []byte, target string, anchor string, tags []goarTypes.Tag) (goarTypes.BundleItem, error)
    GetAddress() string
}
```

## Constants

The package provides default constants for the AO testnet:

```go
const (
    DefaultCuUrl      = "http://cu.ao-testnet.xyz"     // Compute Unit URL
    DefaultMuUrl      = "http://mu.ao-testnet.xyz"     // Message Unit URL
    DefaultAuthority  = "fcoN_xJeisVsPXA-trzVAuIiqO3ydLQxM-L4XbrQKzY"
    DefaultScheduler  = "_GQ33BkPtZrqxA84vM8Zk-N2aO0toNNu_C-l-rawrBA"
    AOSModule        = "JArYBF-D8q2OmZ4Mok00sD2Y_6SYEQ7Hjx-6VZ_jl3g"
    StubArweaveId    = "JArYBF-D8q2OmZ4Mok00sD2Y_6SYEQ7Hjx-6VZ_jl3g"
)
```

## Examples

### Working with Existing Process

```go
// Connect to existing process
existingProcessId := "your-existing-process-id"
client := ao.NewProcessClient(&existingProcessId, ao.DefaultCuUrl, ao.DefaultMuUrl, signer)

// Read current state
result, err := client.Read(ao.DryRunInput{
    Tags: []goarTypes.Tag{
        {Name: "Action", Value: "Info"},
    },
})
```

## Error Handling

The package returns standard Go errors. Common error scenarios:

- **Invalid wallet/signer**: Check wallet format and permissions
- **Process not found**: Verify process ID exists on the network
- **Network errors**: Check CU/MU URLs and network connectivity
- **Invalid Lua code**: Verify syntax when sending code to AOS processes

```go
_, result, err := client.Write(ao.WriteInput{Data: "return 'test'"})
if err != nil {
    fmt.Printf("Write failed: %v\n", err)
    return
}

if result.Error != nil {
    fmt.Printf("Process error: %v\n", result.Error)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

