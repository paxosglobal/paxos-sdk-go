// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// TaxDocumentationService provides access to the TaxDocumentation API.
type TaxDocumentationService struct {
	client *Client
}

// CreateTaxDocumentation createTaxDocumentation stores a partner-collected tax form for an Identity.
//
// Requires scope: identity:write_identity
func (s *TaxDocumentationService) CreateTaxDocumentation(ctx context.Context, body *CreateTaxDocumentationRequest, opts ...RequestOption) (*TaxDocumentation, error) {
	return postTyped[TaxDocumentation](s.client, ctx, "/identity/tax-documentation", body, opts...)
}

// GetTaxDocumentation getTaxDocumentation returns the current accepted tax form for an Identity.
//
// Requires scope: identity:read_identity
func (s *TaxDocumentationService) GetTaxDocumentation(ctx context.Context, identityID string, opts ...RequestOption) (*TaxDocumentation, error) {
	return getTyped[TaxDocumentation](s.client, ctx, buildPath("/identity/tax-documentation/{identity_id}", "identity_id", identityID), opts...)
}
