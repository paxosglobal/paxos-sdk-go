// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// SettlementService provides access to the Settlement API.
//
// Use the Settlements API to facilitate simultaneous exchange of pre-funded
// assets. Automate a variety of use cases that require two-party approval,
// including net settlement of over-the-counter trades, withdrawal requests for
// tri-party collateral, payment requests, and bilateral settlement for
// marketplace end users.
//
// The simple request-and-approval workflow supports both one-directional and
// bidirectional transactions to allow a user to receive an asset (for example,
// send USDP) or simultaneously send and receive assets (for example, send USD
// and receive BTC and ETH). Only when all parties are in agreement and assets
// are fully funded does the Settlements API allow for change of custody. Upon
// completion, all settled assets are immediately available for trading,
// transferring, withdrawal or other supported activities.
//
// The Source Profile initiates the transaction for the Target Profile to
// approve. The `DIRECTION` of each asset is always relative to the Source
// Profile and only the owner of the Source Profile can cancel a transaction if
// it hasn't been approved by the Target Profile.
type SettlementService struct {
	client *Client
}

// ListTransactions list Transactions
//
// Requires scope: settlement:read_transaction
func (s *SettlementService) ListTransactions(ctx context.Context, params *ListTransactionsParams, opts ...RequestOption) *Iterator[Transaction] {
	return paginateTyped[Transaction](s.client, ctx, "/settlement/transactions", "items", params, opts...)
}

// CreateTransaction create Transaction
//
// Requires scope: settlement:write_transaction
func (s *SettlementService) CreateTransaction(ctx context.Context, body *CreateTransactionRequest, opts ...RequestOption) (*Transaction, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[Transaction](s.client, ctx, "/settlement/transactions", body, opts...)
}

// GetTransaction get Transaction
//
// Requires scope: settlement:read_transaction
func (s *SettlementService) GetTransaction(ctx context.Context, transactionID string, opts ...RequestOption) (*Transaction, error) {
	return getTyped[Transaction](s.client, ctx, buildPath("/settlement/transactions/{transaction_id}", "transaction_id", transactionID), opts...)
}

// CancelTransaction cancel Transaction
//
// Requires scope: settlement:write_transaction
func (s *SettlementService) CancelTransaction(ctx context.Context, transactionID string, opts ...RequestOption) (*Transaction, error) {
	return deleteTyped[Transaction](s.client, ctx, buildPath("/settlement/transactions/{transaction_id}", "transaction_id", transactionID), opts...)
}

// AffirmTransaction affirm Transaction
//
// Requires scope: settlement:write_transaction
func (s *SettlementService) AffirmTransaction(ctx context.Context, transactionID string, opts ...RequestOption) (*Transaction, error) {
	return putTyped[Transaction](s.client, ctx, buildPath("/settlement/transactions/{transaction_id}/affirm", "transaction_id", transactionID), nil, opts...)
}
