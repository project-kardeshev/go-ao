package ao

import (
	"fmt"
	"os"
	"testing"

	goarTypes "github.com/everFinance/goar/types"
	"github.com/project-kardeshev/go-ao/signers"
)

func TestSpawnProcess(t *testing.T) {
	fmt.Println("SpawnProcess")
	fmt.Println("Loading wallet")
	wallet, err := os.ReadFile("test_wallet.json")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Wallet loaded")

	signer, err := signers.NewArweaveSigner(wallet)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Signer Created")
	processClient := NewProcessClient(
		nil, 
		DefaultCuUrl, 
		DefaultMuUrl, 
		signer,
	)
	fmt.Println("ProcessClient Created")

	fmt.Println("Spawn Process")
	processId, result, err := processClient.Spawn(SpawnInput{
		Module: AOSModule,
		Authority: DefaultAuthority,
		Scheduler: DefaultScheduler,
		Tags: []goarTypes.Tag{
			{Name: "SDK", Value: "github.com/project-kardeshev/go-ao"},
		},
	})

	fmt.Println("Process ID:", processId, "Result:", result, "Error:", err)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Process Spawned")
	fmt.Println("ProcessId:", processId)
	fmt.Println("Result:", result)

	
	
}

func TestReadProcess(t *testing.T) {
	fmt.Println("TestReadProcess")
	fmt.Println("Loading wallet")
	wallet, err := os.ReadFile("test_wallet.json")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Wallet loaded")

	signer, err := signers.NewArweaveSigner(wallet)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Signer Created")

	processClient := NewProcessClient(
		"JASO-FO1NKxAs2HXFtb7SA2e3u1Q2TbOCq-3wM03bVY", // previously spawned process 
		DefaultCuUrl, 
		DefaultMuUrl, 
		signer,
	)
	fmt.Println("ProcessClient Created")

	fmt.Println("Reading process state")
	result, err := processClient.Read(DryRunInput{
		Id:     "test-read",
		Owner:  signer.GetAddress(),  // Add missing field
		From:   signer.GetAddress(),  // Add missing field
		Anchor: "0",
		Data:   "return 'Hello from AO process'",
		Tags: []goarTypes.Tag{
			{Name: "Action", Value: "Read"},
			{Name: "SDK", Value: "github.com/project-kardeshev/go-ao"},
		},
	})

	fmt.Printf("Read result: %+v, Error: %v\n", result, err)

	if err != nil {
		t.Fatal("Failed to read process:", err)
	}

	if result == nil {
		t.Fatal("Read result is nil")
	}

	fmt.Println("Read test completed successfully")
}