// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// CryptoWithdrawalsService provides access to the CryptoWithdrawals API.
type CryptoWithdrawalsService struct {
	client *Client
}

// PutCryptoDestinationAddress put Crypto Destination Addresses
//
// Requires scope: transfer:write_crypto_destination_address
func (s *CryptoWithdrawalsService) PutCryptoDestinationAddress(ctx context.Context, body *PutCryptoDestinationAddressRequest, opts ...RequestOption) (*PutCryptoDestinationAddressResponse, error) {
	return putTyped[PutCryptoDestinationAddressResponse](s.client, ctx, "/transfer/crypto-destination-address", body, opts...)
}

// ListCryptoDestinationAddresses list Crypto Destination Addresses
//
// Requires scope: transfer:read_crypto_destination_address
func (s *CryptoWithdrawalsService) ListCryptoDestinationAddresses(ctx context.Context, params *ListCryptoDestinationAddressesParams, opts ...RequestOption) *Iterator[CryptoDestinationAddress] {
	return paginateTyped[CryptoDestinationAddress](s.client, ctx, "/transfer/crypto-destination-addresses", "addresses", params, opts...)
}

// CreateCryptoWithdrawal create Crypto Withdrawal
//
// Requires scope: transfer:write_crypto_withdrawal
func (s *CryptoWithdrawalsService) CreateCryptoWithdrawal(ctx context.Context, body *CreateCryptoWithdrawalRequest, opts ...RequestOption) (*Transfer, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[Transfer](s.client, ctx, "/transfer/crypto-withdrawals", body, opts...)
}
