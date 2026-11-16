// Quickstart example: Initialize a Paxos client and make a single API call.
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
		paxos.WithScopes(paxos.ScopeEventsReadEvent),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Fetch a single event to verify connectivity.
	event, err := client.Events.GetEvent(context.Background(), "evt_example", nil)
	if err != nil {
		log.Fatalf("failed to get event: %v", err)
	}

	fmt.Printf("Event ID: %s, Type: %s\n", *event.ID, *event.Type)
}
