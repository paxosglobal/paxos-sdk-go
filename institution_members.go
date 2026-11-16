// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// InstitutionMembersService provides access to the InstitutionMembers API.
//
// Institution Members allow you to associate persons with an institution on
// our Platform. Each institution can have multiple members, with each member
// representing the relationship to a person identity.
//
// When creating an institution identity, you must designate at least one
// institution member. This requirement ensures we have complete information
// about the institution's ownership structure and key personnel, which is
// necessary to onboard institutions.
type InstitutionMembersService struct {
	client *Client
}

// AddInstitutionMembers add Institution Members
//
// Requires scope: identity:write_identity
func (s *InstitutionMembersService) AddInstitutionMembers(ctx context.Context, body *AddInstitutionMembersRequest, opts ...RequestOption) (*AddInstitutionMembersResponse, error) {
	return postTyped[AddInstitutionMembersResponse](s.client, ctx, "/identity/institution-members", body, opts...)
}

// DeleteInstitutionMember remove Institution Member
//
// Requires scope: identity:write_identity
func (s *InstitutionMembersService) DeleteInstitutionMember(ctx context.Context, id string, opts ...RequestOption) (*DeleteInstitutionMemberResponse, error) {
	return deleteTyped[DeleteInstitutionMemberResponse](s.client, ctx, buildPath("/identity/institution-members/{id}", "id", id), opts...)
}
