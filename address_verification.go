// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// AddressVerificationService provides access to the AddressVerification API.
type AddressVerificationService struct {
	client *Client
}

// ListVerifications list address verifications
//
// Requires scope: address:read_ownership_proofs
func (s *AddressVerificationService) ListVerifications(ctx context.Context, params *ListVerificationsParams, opts ...RequestOption) *Iterator[Verification] {
	return paginateTyped[Verification](s.client, ctx, "/addresses/verifications", "items", params, opts...)
}

// CreateVerification create address verification
//
// Requires scope: address:write_ownership_proofs
func (s *AddressVerificationService) CreateVerification(ctx context.Context, body *CreateVerificationRequest, opts ...RequestOption) (*Verification, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[Verification](s.client, ctx, "/addresses/verifications", body, opts...)
}

// SubmitSignature submit address verification signature
//
// Requires scope: address:write_ownership_proofs
func (s *AddressVerificationService) SubmitSignature(ctx context.Context, body *SubmitSignatureRequest, opts ...RequestOption) (*Verification, error) {
	return postTyped[Verification](s.client, ctx, "/addresses/verifications/signature", body, opts...)
}

// GetVerification get address verification
//
// Requires scope: address:read_ownership_proofs
func (s *AddressVerificationService) GetVerification(ctx context.Context, id string, opts ...RequestOption) (*Verification, error) {
	return getTyped[Verification](s.client, ctx, buildPath("/addresses/verifications/{id}", "id", id), opts...)
}

// CancelVerification cancel address verification
//
// Requires scope: address:write_ownership_proofs
func (s *AddressVerificationService) CancelVerification(ctx context.Context, id string, opts ...RequestOption) (*Verification, error) {
	return postTyped[Verification](s.client, ctx, buildPath("/addresses/verifications/{id}/cancel", "id", id), nil, opts...)
}
