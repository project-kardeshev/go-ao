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