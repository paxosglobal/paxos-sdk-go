// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// SandboxDepositsService provides access to the SandboxDeposits API.
type SandboxDepositsService struct {
	client *Client
}

// CreateSandboxDeposit create Sandbox Deposit
//
// Requires scope: funding:read_profile
func (s *SandboxDepositsService) CreateSandboxDeposit(ctx context.Context, profileID string, body *CreateSandboxDepositRequest, opts ...RequestOption) (*CreateSandboxDepositResponse, error) {
	return postTyped[CreateSandboxDepositResponse](s.client, ctx, buildPath("/sandbox/profiles/{profile_id}/deposit", "profile_id", profileID), body, opts...)
}
