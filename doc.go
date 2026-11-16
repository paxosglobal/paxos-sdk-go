// Package paxos is the official Go client for the Paxos API.
//
// It provides typed models and one service client per API area, with OAuth2
// client-credentials authentication, retries, typed errors, auto-paginating
// iterators, request signing and a Decimal type for monetary amounts. The SDK
// is in beta, and breaking changes may occur before 1.0.
//
// Create one Client per application component and share it across goroutines.
// Every service is a field on the client, and every method takes a
// context.Context first:
//
//	client, err := paxos.NewClient(
//		os.Getenv("PAXOS_CLIENT_ID"),
//		os.Getenv("PAXOS_CLIENT_SECRET"),
//		paxos.WithSandbox(),
//		paxos.WithScopes(paxos.ScopeTransferReadTransfer),
//	)
//	if err != nil {
//		return err
//	}
//	transfer, err := client.Transfers.GetTransfer(ctx, "txn_123")
//
// List methods return an *Iterator that fetches pages on demand:
//
//	iter := client.Transfers.ListTransfers(ctx, &paxos.ListTransfersParams{Limit: paxos.Int32(50)})
//	for iter.Next() {
//		fmt.Println(iter.Current().ID)
//	}
//	if err := iter.Err(); err != nil {
//		return err
//	}
//
// WithSandbox selects the sandbox environment; omit it for production.
// Sandbox and production credentials differ.
//
// See https://docs.paxos.com/sdks for the SDK guides and
// https://docs.paxos.com/api-reference/introduction for the API reference.
package paxos
