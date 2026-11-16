// Sandbox deposit example: Create a test deposit in the sandbox environment.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	paxos "github.com/paxosglobal/paxos-sdk-go"
)

func main() {
	// Create a client configured for sandbox mode.
	client, err := paxos.NewClient(
		os.Getenv("PAXOS_CLIENT_ID"),
		os.Getenv("PAXOS_CLIENT_SECRET"),
		paxos.WithSandbox(),
		paxos.WithScopes(paxos.ScopeFundingReadProfile),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	cryptoNetwork := paxos.CryptoNetworkETHEREUM

	// Simulate a crypto deposit of 1.5 ETH into a profile.
	deposit, err := client.SandboxDeposits.CreateSandboxDeposit(context.Background(),
		os.Getenv("PAXOS_PROFILE_ID"),
		&paxos.CreateSandboxDepositRequest{
			Asset:         "ETH",
			Amount:        paxos.MustDecimal("1.5"),
			CryptoNetwork: &cryptoNetwork,
		},
	)
	if err != nil {
		log.Fatalf("failed to create sandbox deposit: %v", err)
	}

	fmt.Printf("Sandbox deposit created: activity_id=%s\n", *deposit.ActivityID)
}
