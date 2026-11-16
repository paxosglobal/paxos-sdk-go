// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// FiatTransfersService provides access to the FiatTransfers API.
//
// The Fiat Transfers API allows clients to deposit and withdraw fiat via the
// following Fiat Networks: Wire (Fedwire and SWIFT) and CUBIX. The supported
// asset is USD.
//
// The full set of Fiat Transfers API endpoints is:
//
//   - For Deposits:
//   - [FiatTransfersService.CreateFiatDepositInstructions]
//   - [FiatTransfersService.ListFiatDepositInstructions]
//   - [FiatTransfersService.GetFiatDepositInstructions]
//   - For Withdrawals:
//   - [FiatTransfersService.CreateFiatAccount]
//   - [FiatTransfersService.ListFiatAccounts]
//   - [FiatTransfersService.GetFiatAccount]
//   - [FiatTransfersService.DeleteFiatAccount]
//   - [FiatTransfersService.UpdateFiatAccount]
//   - [FiatTransfersService.CreateFiatWithdrawal]
type FiatTransfersService struct {
	client *Client
}

// ListFiatAccounts list Fiat Accounts
//
// Requires scope: transfer:read_fiat_account
func (s *FiatTransfersService) ListFiatAccounts(ctx context.Context, params *ListFiatAccountsParams, opts ...RequestOption) *Iterator[FiatAccount] {
	return paginateTyped[FiatAccount](s.client, ctx, "/transfer/fiat-accounts", "items", params, opts...)
}

// CreateFiatAccount create Fiat Account
//
// Requires scope: transfer:write_fiat_account
func (s *FiatTransfersService) CreateFiatAccount(ctx context.Context, body *CreateFiatAccountRequest, opts ...RequestOption) (*FiatAccount, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[FiatAccount](s.client, ctx, "/transfer/fiat-accounts", body, opts...)
}

// GetFiatAccount get Fiat Account
//
// Requires scope: transfer:read_fiat_account
func (s *FiatTransfersService) GetFiatAccount(ctx context.Context, id string, params *GetFiatAccountParams, opts ...RequestOption) (*FiatAccount, error) {
	return getWithParamsTyped[FiatAccount](s.client, ctx, buildPath("/transfer/fiat-accounts/{id}", "id", id), params, opts...)
}

// DeleteFiatAccount delete Fiat Account
//
// Requires scope: transfer:write_fiat_account
func (s *FiatTransfersService) DeleteFiatAccount(ctx context.Context, id string, opts ...RequestOption) (*DeleteFiatAccountResponse, error) {
	return deleteTyped[DeleteFiatAccountResponse](s.client, ctx, buildPath("/transfer/fiat-accounts/{id}", "id", id), opts...)
}

// UpdateFiatAccount update Fiat Account
//
// Requires scope: transfer:write_fiat_account
func (s *FiatTransfersService) UpdateFiatAccount(ctx context.Context, id string, body *UpdateFiatAccountRequest, opts ...RequestOption) (*FiatAccount, error) {
	return putTyped[FiatAccount](s.client, ctx, buildPath("/transfer/fiat-accounts/{id}", "id", id), body, opts...)
}

// ListFiatDepositInstructions list Fiat Deposit Instructions
//
// Requires scope: transfer:read_fiat_deposit_instructions
func (s *FiatTransfersService) ListFiatDepositInstructions(ctx context.Context, params *ListFiatDepositInstructionsParams, opts ...RequestOption) *Iterator[FiatDepositInstructions] {
	return paginateTyped[FiatDepositInstructions](s.client, ctx, "/transfer/fiat-deposit-instructions", "items", params, opts...)
}

// CreateFiatDepositInstructions create Fiat Deposit Instructions
//
// Requires scope: transfer:write_fiat_deposit_instructions
func (s *FiatTransfersService) CreateFiatDepositInstructions(ctx context.Context, body *CreateFiatDepositInstructionsRequest, opts ...RequestOption) (*FiatDepositInstructions, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[FiatDepositInstructions](s.client, ctx, "/transfer/fiat-deposit-instructions", body, opts...)
}

// GetFiatDepositInstructions get Fiat Deposit Instructions
//
// Requires scope: transfer:read_fiat_deposit_instructions
func (s *FiatTransfersService) GetFiatDepositInstructions(ctx context.Context, id string, opts ...RequestOption) (*FiatDepositInstructions, error) {
	return getTyped[FiatDepositInstructions](s.client, ctx, buildPath("/transfer/fiat-deposit-instructions/{id}", "id", id), opts...)
}

// CreateFiatWithdrawal create Fiat Withdrawal
//
// Requires scope: transfer:write_fiat_withdrawal
func (s *FiatTransfersService) CreateFiatWithdrawal(ctx context.Context, body *CreateFiatWithdrawalRequest, opts ...RequestOption) (*Transfer, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[Transfer](s.client, ctx, "/transfer/fiat-withdrawals", body, opts...)
}
