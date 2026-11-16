// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// AccountsService provides access to the Accounts API.
//
// Accounts connect Identities (persons or institutions) to Profiles (asset
// balances). Key features include:
//
//   - Balance Management: Track and manage profile balances via Subledgering
//   - Joint Account Support: Link multiple Identities to a single Account for
//     shared on platform balance
//
// Accounts are required for all Identity API integrations.
type AccountsService struct {
	client *Client
}

// ListAccounts list Accounts
//
// Requires scope: identity:read_account
func (s *AccountsService) ListAccounts(ctx context.Context, params *ListAccountsParams, opts ...RequestOption) *Iterator[Account] {
	return paginateTyped[Account](s.client, ctx, "/identity/accounts", "items", params, opts...)
}

// CreateAccount create Account
//
// Requires scope: identity:write_account
func (s *AccountsService) CreateAccount(ctx context.Context, body *CreateAccountRequest, opts ...RequestOption) (*Account, error) {
	if body != nil && (body.Account.RefID == nil || *body.Account.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.Account.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[Account](s.client, ctx, "/identity/accounts", body, opts...)
}

// UpdateAccount update Account
//
// Requires scope: identity:write_account
func (s *AccountsService) UpdateAccount(ctx context.Context, body *UpdateAccountRequest, opts ...RequestOption) (*Account, error) {
	return putTyped[Account](s.client, ctx, "/identity/accounts", body, opts...)
}

// GetAccount calls the GetAccount API operation.
//
// Requires scope: identity:read_account
func (s *AccountsService) GetAccount(ctx context.Context, id string, opts ...RequestOption) (*Account, error) {
	return getTyped[Account](s.client, ctx, buildPath("/identity/accounts/{id}", "id", id), opts...)
}
