// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// SandboxFiatTransfersService provides access to the SandboxFiatTransfers API.
type SandboxFiatTransfersService struct {
	client *Client
}

// InitiateSandboxFiatDeposit initiate Sandbox Fiat Deposit
//
// Requires scope: transfer:write_sandbox_fiat_deposit
func (s *SandboxFiatTransfersService) InitiateSandboxFiatDeposit(ctx context.Context, body *InitiateSandboxFiatDepositRequest, opts ...RequestOption) (*InitiateSandboxFiatDepositResponse, error) {
	return postTyped[InitiateSandboxFiatDepositResponse](s.client, ctx, "/sandbox/fiat-deposits", body, opts...)
}
