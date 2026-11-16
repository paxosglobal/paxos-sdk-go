// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// CryptoDepositsService provides access to the CryptoDeposits API.
//
// For Paxos Global PTE Ltd users, update or reject a crypto deposit that does
// not contain the required travel rule information. Deposits with the
// `NEEDS_INFO` status require additional sender information before they can be
// processed.
//
//   - For individuals, both first and last name are required.
//   - For institutions, the institution name is required.
//
// If a deposit is rejected, the funds will no longer be available.
type CryptoDepositsService struct {
	client *Client
}

// RejectCryptoDeposit reject Crypto Deposit
//
// Requires scope: transfer:reject_crypto_deposit
func (s *CryptoDepositsService) RejectCryptoDeposit(ctx context.Context, id string, body *RejectCryptoDepositRequest, opts ...RequestOption) (*RejectCryptoDepositResponse, error) {
	return postTyped[RejectCryptoDepositResponse](s.client, ctx, buildPath("/transfer/crypto-deposits/{id}/reject", "id", id), body, opts...)
}

// UpdateCryptoDeposit update Crypto Deposit
//
// Requires scope: transfer:update_crypto_deposit
func (s *CryptoDepositsService) UpdateCryptoDeposit(ctx context.Context, id string, body *UpdateCryptoDepositRequest, opts ...RequestOption) (*UpdateCryptoDepositResponse, error) {
	return postTyped[UpdateCryptoDepositResponse](s.client, ctx, buildPath("/transfer/crypto-deposits/{id}/update", "id", id), body, opts...)
}
