// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// FeesService provides access to the Fees API.
type FeesService struct {
	client *Client
}

// CreateCryptoWithdrawalFee create Crypto Withdrawal Fee
//
// Requires scope: fee:write_crypto_withdrawal_fee
func (s *FeesService) CreateCryptoWithdrawalFee(ctx context.Context, body *CreateCryptoWithdrawalFeeRequest, opts ...RequestOption) (*CryptoWithdrawalFee, error) {
	return postTyped[CryptoWithdrawalFee](s.client, ctx, "/transfer/fees/crypto-withdrawal", body, opts...)
}
