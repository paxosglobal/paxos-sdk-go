// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// SandboxIdentityService provides access to the SandboxIdentity API.
type SandboxIdentityService struct {
	client *Client
}

// SandboxBackdateIdentityClosure sandbox Backdate Identity Closure
//
// Requires scope: identity:write_identity
func (s *SandboxIdentityService) SandboxBackdateIdentityClosure(ctx context.Context, identityID string, body *SandboxBackdateIdentityClosureRequest, opts ...RequestOption) (*SandboxBackdateIdentityClosureResponse, error) {
	return postTyped[SandboxBackdateIdentityClosureResponse](s.client, ctx, buildPath("/identity/identities/{identity_id}/backdate-closure", "identity_id", identityID), body, opts...)
}

// SandboxStartKycRefresh sandbox Start KYC Refresh
//
// Requires scope: identity:write_identity
func (s *SandboxIdentityService) SandboxStartKycRefresh(ctx context.Context, identityID string, opts ...RequestOption) (*SandboxStartKycRefreshResponse, error) {
	return postTyped[SandboxStartKycRefreshResponse](s.client, ctx, buildPath("/identity/identities/{identity_id}/kyc-refresh", "identity_id", identityID), nil, opts...)
}

// SandboxSetIdentityStatus sandbox Set Identity Status
//
// Requires scope: identity:write_identity
func (s *SandboxIdentityService) SandboxSetIdentityStatus(ctx context.Context, id string, body *SandboxSetIdentityStatusRequest, opts ...RequestOption) (*SandboxSetIdentityResponse, error) {
	return putTyped[SandboxSetIdentityResponse](s.client, ctx, buildPath("/identity/identities/{id}/sandbox-status", "id", id), body, opts...)
}

// SandboxForceStartKycRefresh sandbox Force Start KYC Refresh
//
// Requires scope: identity:write_identity
func (s *SandboxIdentityService) SandboxForceStartKycRefresh(ctx context.Context, body *SandboxForceStartKycRefreshRequest, opts ...RequestOption) (*SandboxForceStartKycRefreshResponse, error) {
	return postTyped[SandboxForceStartKycRefreshResponse](s.client, ctx, "/identity/kycrefresh/force-start", body, opts...)
}

// SandboxResetKycRefresh sandbox Reset KYC Refresh
//
// Requires scope: identity:write_identity
func (s *SandboxIdentityService) SandboxResetKycRefresh(ctx context.Context, body *SandboxResetKycRefreshRequest, opts ...RequestOption) (*SandboxResetKycRefreshResponse, error) {
	return postTyped[SandboxResetKycRefreshResponse](s.client, ctx, "/identity/kycrefresh/reset", body, opts...)
}
