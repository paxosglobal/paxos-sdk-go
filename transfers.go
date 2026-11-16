// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// TransfersService provides access to the Transfers API.
//
// Each deposit to and withdrawal from a Paxos account is a Transfer. A
// Transfer is associated with a particular Profile, and increases or decreases
// the available balance of a single asset in that Profile.
type TransfersService struct {
	client *Client
}

// ListTransfers list Transfers
//
// Requires scope: transfer:read_transfer
func (s *TransfersService) ListTransfers(ctx context.Context, params *ListTransfersParams, opts ...RequestOption) *Iterator[Transfer] {
	return paginateTyped[Transfer](s.client, ctx, "/transfer/transfers", "items", params, opts...)
}

// GetTransfer get Transfer
//
// Requires scope: transfer:read_transfer
func (s *TransfersService) GetTransfer(ctx context.Context, id string, opts ...RequestOption) (*Transfer, error) {
	return getTyped[Transfer](s.client, ctx, buildPath("/transfer/transfers/{id}", "id", id), opts...)
}
