// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// AccountMembersService provides access to the AccountMembers API.
//
// Account Members allow you to associate an Identity with an Account on our
// Platform. Each Account can be linked to multiple Identities, which supports
// Joint Accounts (where two Identities have the ability to transact against a
// single Profile)
//
// Account members are effectively immutable. In order to update an account
// member, clients should remove then re-add the account member with the
// desired state.
//
// The API requires clients to have write access on the account and all
// identities associated with the account.
type AccountMembersService struct {
	client *Client
}

// AddAccountMembers add Account Members
//
// Requires scope: identity:write_account
func (s *AccountMembersService) AddAccountMembers(ctx context.Context, body *AddAccountMembersRequest, opts ...RequestOption) (*AddAccountMembersResponse, error) {
	return postTyped[AddAccountMembersResponse](s.client, ctx, "/identity/account-members", body, opts...)
}

// DeleteAccountMember remove Account Member
//
// Requires scope: identity:write_account
func (s *AccountMembersService) DeleteAccountMember(ctx context.Context, id string, opts ...RequestOption) (*DeleteAccountMemberResponse, error) {
	return deleteTyped[DeleteAccountMemberResponse](s.client, ctx, buildPath("/identity/account-members/{id}", "id", id), opts...)
}
