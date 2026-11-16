// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// DepositAddressesService provides access to the DepositAddresses API.
type DepositAddressesService struct {
	client *Client
}

// ListDepositAddresses list Deposit Addresses
//
// Requires scope: transfer:read_deposit_address
func (s *DepositAddressesService) ListDepositAddresses(ctx context.Context, params *ListDepositAddressesParams, opts ...RequestOption) *Iterator[DepositAddress] {
	return paginateTyped[DepositAddress](s.client, ctx, "/transfer/deposit-addresses", "items", params, opts...)
}

// CreateDepositAddress create Deposit Address
//
// Requires scope: transfer:write_deposit_address
func (s *DepositAddressesService) CreateDepositAddress(ctx context.Context, body *CreateDepositAddressRequest, opts ...RequestOption) (*DepositAddress, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[DepositAddress](s.client, ctx, "/transfer/deposit-addresses", body, opts...)
}
