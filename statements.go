// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// StatementsService provides access to the Statements API.
//
// A Statement summarizes all eligible stablecoin rewards for a given
// organization over a fixed time period. It aggregates earned reward types
// (e.g., custody, mint, acceptance) and records whether the statement was
// paid.
//
// Statements are generated monthly and are immutable. Revised statements will
// be issued as a new statement with a unique ID. Statements can be used to
// reconcile balances or for accounting and accounts receivable purposes.
//
// *These endpoints are only relevant for Global Dollar Network (GDN)
// partners.*
type StatementsService struct {
	client *Client
}

// ListStatements list Statements
//
// Requires scope: statements:read_statement
//
// Deprecated: This API operation is deprecated.
func (s *StatementsService) ListStatements(ctx context.Context, params *ListStatementsParams, opts ...RequestOption) *Iterator[Statement] {
	return paginateTyped[Statement](s.client, ctx, "/statements", "statements", params, opts...)
}
