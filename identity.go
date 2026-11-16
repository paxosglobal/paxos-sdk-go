// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// IdentityService provides access to the Identity API.
//
// An Identity represents a person or institution who is able to take actions
// on the Paxos Platform. Depending on your integration type you might not need
// to create Identities for your end users.
//
// Identities must pass a set of verifications and checks before they're able
// to transact on the Platform. Learn more about getting started with Identity
// and managing their onboarding lifecycle in the Identity Developer Guide.
type IdentityService struct {
	client *Client
}

// ListIdentities list Identities
//
// Requires scope: identity:read_identity
func (s *IdentityService) ListIdentities(ctx context.Context, params *ListIdentitiesParams, opts ...RequestOption) *Iterator[Identity] {
	return paginateTyped[Identity](s.client, ctx, "/identity/identities", "items", params, opts...)
}

// CreateIdentity create Identity
//
// Requires scope: identity:write_identity
func (s *IdentityService) CreateIdentity(ctx context.Context, body *CreateIdentityRequest, opts ...RequestOption) (*Identity, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[Identity](s.client, ctx, "/identity/identities", body, opts...)
}

// GetIdentity get Identity
//
// Requires scope: identity:read_identity
func (s *IdentityService) GetIdentity(ctx context.Context, id string, params *GetIdentityParams, opts ...RequestOption) (*Identity, error) {
	return getWithParamsTyped[Identity](s.client, ctx, buildPath("/identity/identities/{id}", "id", id), params, opts...)
}

// UpdateIdentity update Identity
//
// Requires scope: identity:write_identity
func (s *IdentityService) UpdateIdentity(ctx context.Context, id string, body *UpdateIdentityRequest, opts ...RequestOption) (*Identity, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	return putTyped[Identity](s.client, ctx, buildPath("/identity/identities/{id}", "id", id), body, opts...)
}
