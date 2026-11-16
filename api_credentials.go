// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// APICredentialsService provides access to the APICredentials API.
//
// API credentials allow programmatic access to Paxos APIs. These endpoints
// enable you to manage and view your API credentials.
type APICredentialsService struct {
	client *Client
}

// ListAPICredentials list API Credentials
//
// Requires scope: api_creds:read_credentials
func (s *APICredentialsService) ListAPICredentials(ctx context.Context, params *ListAPICredentialsParams, opts ...RequestOption) *Iterator[FirstPartyClient] {
	return paginateTyped[FirstPartyClient](s.client, ctx, "/api-creds/credentials", "items", params, opts...)
}

// GetAPICredentialsSelf calls the GetAPICredentialsSelf API operation.
func (s *APICredentialsService) GetAPICredentialsSelf(ctx context.Context, opts ...RequestOption) (*FirstPartyClient, error) {
	return getTyped[FirstPartyClient](s.client, ctx, "/api-creds/credentials/self", opts...)
}

// DeleteAPICredentials delete API Credentials
//
// Requires scope: api_creds:delete_credentials
func (s *APICredentialsService) DeleteAPICredentials(ctx context.Context, clientID string, opts ...RequestOption) (*DeleteAPICredentialsResponse, error) {
	return deleteTyped[DeleteAPICredentialsResponse](s.client, ctx, buildPath("/api-creds/credentials/{client_id}", "client_id", clientID), opts...)
}
