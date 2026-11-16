// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// IdentityDocumentsService provides access to the IdentityDocuments API.
//
// The Identity Documents API allows clients to send documents to Paxos to
// validate identities during the KYC process. Files can be updated multiple
// times by re-sending the document type for processing.
//
// We currently accept the following file formats:
//
//   - `pdf`
//   - `jpg`
//   - `png`
//
// Steps to upload a document are as follows:
//
//  1. Send a `PUT` request to `/identity/identities/{id}/documents` to get
//     the upload URL.
//  2. Send your document to the provided URL via a PUT request. Take the
//     following example using curl: `$ curl --upload-file
//     Proof_of_residence.jpg $URI `
//
// Uploaded document metadata will be available via `GET` request to
// `/identity/identities/{id}/documents`
type IdentityDocumentsService struct {
	client *Client
}

// ListIdentityDocuments list Identity Documents
//
// Requires scope: identity:read_identity
func (s *IdentityDocumentsService) ListIdentityDocuments(ctx context.Context, identityID string, params *ListIdentityDocumentsParams, opts ...RequestOption) (*ListIdentityDocumentsResponse, error) {
	return getWithParamsTyped[ListIdentityDocumentsResponse](s.client, ctx, buildPath("/identity/identities/{identity_id}/documents", "identity_id", identityID), params, opts...)
}

// DocumentUpload document Upload
//
// Requires scope: identity:write_identity
func (s *IdentityDocumentsService) DocumentUpload(ctx context.Context, identityID string, body *DocumentUploadRequest, opts ...RequestOption) (*DocumentUploadResponse, error) {
	return putTyped[DocumentUploadResponse](s.client, ctx, buildPath("/identity/identities/{identity_id}/documents", "identity_id", identityID), body, opts...)
}
