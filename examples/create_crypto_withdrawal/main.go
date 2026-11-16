// Create crypto withdrawal example: End-to-end transfer with ref_id for idempotency.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	paxos "github.com/paxosglobal/paxos-sdk-go"
)

func main() {
	// Create a client using environment variables for credentials.
	client, err := paxos.NewClient(
		os.Getenv("PAXOS_CLIENT_ID"),
		os.Getenv("PAXOS_CLIENT_SECRET"),
		paxos.WithSandbox(),
		paxos.WithScopes(paxos.ScopeTransferWriteCryptoWithdrawal),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	amount := paxos.MustDecimal("0.001")

	// Create a crypto withdrawal with a unique ref_id for idempotency.
	withdrawal, err := client.CryptoWithdrawals.CreateCryptoWithdrawal(context.Background(),
		&paxos.CreateCryptoWithdrawalRequest{
			ProfileID:          os.Getenv("PAXOS_PROFILE_ID"),
			Asset:              "BTC",
			Amount:             &amount,
			DestinationAddress: "bc1qexampleaddress",
			RefID:              paxos.String("withdrawal-example-001"),
			CryptoNetwork:      paxos.CryptoNetworkBITCOIN,
		},
	)
	if err != nil {
		log.Fatalf("failed to create crypto withdrawal: %v", err)
	}

	fmt.Printf("Withdrawal ID: %s\n", withdrawal.ID)
	fmt.Printf("Status: %s\n", withdrawal.Status)
	fmt.Printf("Amount: %s %s\n", withdrawal.Amount, withdrawal.Asset)
}
