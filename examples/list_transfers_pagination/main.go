// List transfers with pagination: Demonstrate auto-pagination iterator.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	paxos "github.com/paxosglobal/paxos-sdk-go"
)

func main() {
	// Create a client using environment variables for credentials.
	client, err := paxos.NewClient(
		os.Getenv("PAXOS_CLIENT_ID"),
		os.Getenv("PAXOS_CLIENT_SECRET"),
		paxos.WithSandbox(),
		paxos.WithScopes(paxos.ScopeTransferReadTransfer),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// List transfers using the auto-pagination iterator.
	// The iterator automatically fetches additional pages as needed.
	// Limit is a page-size hint (10 items per page here), not a cap on the
	// total — to bound the whole iteration instead, pass WithMaxItems
	// and/or WithMaxDuration.
	iter := client.Transfers.ListTransfers(context.Background(),
		&paxos.ListTransfersParams{
			Limit: paxos.Int32(10),
		},
		paxos.WithMaxItems(25), // stop after 25 transfers total (~3 pages)
		paxos.WithMaxDuration(30*time.Second),
	)

	count := 0
	for iter.Next() {
		transfer := iter.Current()
		fmt.Printf("Transfer %s: %s %s (%s)\n",
			transfer.ID, transfer.Amount, transfer.Asset, transfer.Status)
		count++
	}

	// Always check for errors after iteration completes.
	if err := iter.Err(); err != nil {
		log.Fatalf("error listing transfers: %v", err)
	}

	fmt.Printf("\nTotal transfers: %d\n", count)
}
