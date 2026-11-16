// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

// Scope is an OAuth2 scope accepted by WithScopes.
type Scope string

const (
	// ScopeAddressReadOwnershipProofs: Retrieve address ownership verification details
	ScopeAddressReadOwnershipProofs Scope = "address:read_ownership_proofs"
	// ScopeAddressWriteOwnershipProofs: Create, sign or cancel an address ownership verification
	ScopeAddressWriteOwnershipProofs Scope = "address:write_ownership_proofs"
	// ScopeAPICredsDeleteCredentials: Delete API credentials
	ScopeAPICredsDeleteCredentials Scope = "api_creds:delete_credentials"
	// ScopeAPICredsReadCredentials: Retrieve API credentials details
	ScopeAPICredsReadCredentials Scope = "api_creds:read_credentials"
	// ScopeConversionReadConversionStablecoin: Retrieve stablecoin conversion details
	ScopeConversionReadConversionStablecoin Scope = "conversion:read_conversion_stablecoin"
	// ScopeConversionWriteConversionStablecoin: Create or cancel a stablecoin conversion
	ScopeConversionWriteConversionStablecoin Scope = "conversion:write_conversion_stablecoin"
	// ScopeEventsReadEvent: Retrieve events
	ScopeEventsReadEvent Scope = "events:read_event"
	// ScopeExchangeHistoricalPrices: Retrieve market average prices at a certain time increment
	ScopeExchangeHistoricalPrices Scope = "exchange:historical_prices"
	// ScopeExchangeReadIssuerQuoteExecution: Retrieve issuer quote execution details
	ScopeExchangeReadIssuerQuoteExecution Scope = "exchange:read_issuer_quote_execution"
	// ScopeExchangeReadOrder: Retrieve order or order execution details
	ScopeExchangeReadOrder Scope = "exchange:read_order"
	// ScopeExchangeReadQuote: Retrieve quote details for buying or selling an asset
	ScopeExchangeReadQuote Scope = "exchange:read_quote"
	// ScopeExchangeReadQuoteExecution: Retrieve quote execution details
	ScopeExchangeReadQuoteExecution Scope = "exchange:read_quote_execution"
	// ScopeExchangeWriteIssuerQuote: Create an issuer quote or retrieve issuer quote settlement availability
	ScopeExchangeWriteIssuerQuote Scope = "exchange:write_issuer_quote"
	// ScopeExchangeWriteIssuerQuoteExecution: Create an issuer quote execution
	ScopeExchangeWriteIssuerQuoteExecution Scope = "exchange:write_issuer_quote_execution"
	// ScopeExchangeWriteOrder: Create or cancel an order for buying or selling an asset
	ScopeExchangeWriteOrder Scope = "exchange:write_order"
	// ScopeExchangeWriteQuoteExecution: Create a quote execution for buying or selling an asset
	ScopeExchangeWriteQuoteExecution Scope = "exchange:write_quote_execution"
	// ScopeFeeWriteCryptoWithdrawalFee: Create a guaranteed fee for crypto withdrawal
	ScopeFeeWriteCryptoWithdrawalFee Scope = "fee:write_crypto_withdrawal_fee"
	// ScopeFundingReadBankBalance: Retrieve Paxos dedicated bank account balance
	ScopeFundingReadBankBalance Scope = "funding:read_bank_balance"
	// ScopeFundingReadProfile: Retrieve Profile details and deposit funds in Sandbox
	ScopeFundingReadProfile Scope = "funding:read_profile"
	// ScopeFundingWriteProfile: Create a Profile
	ScopeFundingWriteProfile Scope = "funding:write_profile"
	// ScopeIdentityReadAccount: Retrieve Account details
	ScopeIdentityReadAccount Scope = "identity:read_account"
	// ScopeIdentityReadIdentity: Retrieve Identity details or documents
	ScopeIdentityReadIdentity Scope = "identity:read_identity"
	// ScopeIdentityReadIdentityControl: Retrieve Identity Control details
	ScopeIdentityReadIdentityControl Scope = "identity:read_identity_control"
	// ScopeIdentityWriteAccount: Create or update Account and Account Members
	ScopeIdentityWriteAccount Scope = "identity:write_account"
	// ScopeIdentityWriteIdentity: Create or update Identity details and set Sandbox Identify Status
	ScopeIdentityWriteIdentity Scope = "identity:write_identity"
	// ScopeIdentityWriteIdentityControl: Create or delete an Identity Control
	ScopeIdentityWriteIdentityControl Scope = "identity:write_identity_control"
	// ScopeOrchestrationReadOrchestration: Retrieve Orchestration details
	ScopeOrchestrationReadOrchestration Scope = "orchestration:read_orchestration"
	// ScopeOrchestrationReadOrchestrationRule: Retrieve Orchestration Rule details
	ScopeOrchestrationReadOrchestrationRule Scope = "orchestration:read_orchestration_rule"
	// ScopeOrchestrationWriteOrchestration: Create an Orchestration
	ScopeOrchestrationWriteOrchestration Scope = "orchestration:write_orchestration"
	// ScopeOrchestrationWriteOrchestrationRule: Create or delete an Orchestration Rule
	ScopeOrchestrationWriteOrchestrationRule Scope = "orchestration:write_orchestration_rule"
	// ScopeRewardsCreatePayoutGroups: Create a rewards Payout Group
	ScopeRewardsCreatePayoutGroups Scope = "rewards:create_payout_groups"
	// ScopeRewardsDeletePayoutGroups: Delete a rewards Payout Group
	ScopeRewardsDeletePayoutGroups Scope = "rewards:delete_payout_groups"
	// ScopeRewardsReadClaims: Retrieve rewards claims
	ScopeRewardsReadClaims Scope = "rewards:read_claims"
	// ScopeRewardsReadClaimsSchedule: Retrieve rewards claim schedules
	ScopeRewardsReadClaimsSchedule Scope = "rewards:read_claims_schedule"
	// ScopeRewardsReadMonitoringAddress: Retrieve monitoring address and reward address details
	ScopeRewardsReadMonitoringAddress Scope = "rewards:read_monitoring_address"
	// ScopeRewardsReadPayoutGroups: Retrieve rewards Payout Group details
	ScopeRewardsReadPayoutGroups Scope = "rewards:read_payout_groups"
	// ScopeRewardsReadRewardsDetails: Retrieve rewards details and rates
	ScopeRewardsReadRewardsDetails Scope = "rewards:read_rewards_details"
	// ScopeRewardsReadStatements: Retrieve rewards statements
	ScopeRewardsReadStatements Scope = "rewards:read_statements"
	// ScopeRewardsUpdatePayoutGroups: Update a rewards Payout Group
	ScopeRewardsUpdatePayoutGroups Scope = "rewards:update_payout_groups"
	// ScopeRewardsWriteClaimsSchedule: Create, update or delete a rewards claim schedule
	ScopeRewardsWriteClaimsSchedule Scope = "rewards:write_claims_schedule"
	// ScopeRewardsWriteMonitoringAddress: Create, update or delete a monitoring address or reward address
	ScopeRewardsWriteMonitoringAddress Scope = "rewards:write_monitoring_address"
	// ScopeRewardsWritePayoutGroups: Create, update or delete a rewards Payout Group
	ScopeRewardsWritePayoutGroups Scope = "rewards:write_payout_groups"
	// ScopeSettlementReadTransaction: Retrieve settlement transaction details
	ScopeSettlementReadTransaction Scope = "settlement:read_transaction"
	// ScopeSettlementWriteTransaction: Create, affirm or cancel a settlement transaction
	ScopeSettlementWriteTransaction Scope = "settlement:write_transaction"
	// ScopeStatementsReadPayment: Retrieve payment details
	ScopeStatementsReadPayment Scope = "statements:read_payment"
	// ScopeStatementsReadStatement: Retrieve statement details
	ScopeStatementsReadStatement Scope = "statements:read_statement"
	// ScopeTaxReadTaxForm: Retrieve tax details
	ScopeTaxReadTaxForm Scope = "tax:read_tax_form"
	// ScopeTaxReadTaxLot: Retrieve tax lot details
	ScopeTaxReadTaxLot Scope = "tax:read_tax_lot"
	// ScopeTaxWriteTaxLot: Update the given tax-lot ID
	ScopeTaxWriteTaxLot Scope = "tax:write_tax_lot"
	// ScopeTransferReadCryptoDestinationAddress: Retrieve crypto destination address details
	ScopeTransferReadCryptoDestinationAddress Scope = "transfer:read_crypto_destination_address"
	// ScopeTransferReadDepositAddress: Retrieve deposit address details
	ScopeTransferReadDepositAddress Scope = "transfer:read_deposit_address"
	// ScopeTransferReadFiatAccount: Retrieve Fiat Account details
	ScopeTransferReadFiatAccount Scope = "transfer:read_fiat_account"
	// ScopeTransferReadFiatDepositInstructions: Retrieve fiat deposit instruction details
	ScopeTransferReadFiatDepositInstructions Scope = "transfer:read_fiat_deposit_instructions"
	// ScopeTransferReadTransfer: Retrieve transfer details
	ScopeTransferReadTransfer Scope = "transfer:read_transfer"
	// ScopeTransferReadTransferLimit: Retrieve limits for the given transaction type
	ScopeTransferReadTransferLimit Scope = "transfer:read_transfer_limit"
	// ScopeTransferRejectCryptoDeposit: Reject a crypto deposit (travel rule)
	ScopeTransferRejectCryptoDeposit Scope = "transfer:reject_crypto_deposit"
	// ScopeTransferUpdateCryptoDeposit: Provide required travel-rule details
	ScopeTransferUpdateCryptoDeposit Scope = "transfer:update_crypto_deposit"
	// ScopeTransferWriteCryptoDestinationAddress: Create or update crypto destination addresses
	ScopeTransferWriteCryptoDestinationAddress Scope = "transfer:write_crypto_destination_address"
	// ScopeTransferWriteCryptoWithdrawal: Withdraw asset to a specified destination address
	ScopeTransferWriteCryptoWithdrawal Scope = "transfer:write_crypto_withdrawal"
	// ScopeTransferWriteDepositAddress: Create an deposit address on a blockchain network
	ScopeTransferWriteDepositAddress Scope = "transfer:write_deposit_address"
	// ScopeTransferWriteFiatAccount: Create, update or delete a Fiat Account
	ScopeTransferWriteFiatAccount Scope = "transfer:write_fiat_account"
	// ScopeTransferWriteFiatDepositInstructions: Create, update or delete fiat deposit instructions
	ScopeTransferWriteFiatDepositInstructions Scope = "transfer:write_fiat_deposit_instructions"
	// ScopeTransferWriteFiatWithdrawal: Withdraw fiat to the given destination
	ScopeTransferWriteFiatWithdrawal Scope = "transfer:write_fiat_withdrawal"
	// ScopeTransferWriteInternalTransfer: Transfer assets between two Profiles
	ScopeTransferWriteInternalTransfer Scope = "transfer:write_internal_transfer"
	// ScopeTransferWritePaxosTransfer: Create a Paxos transfer
	ScopeTransferWritePaxosTransfer Scope = "transfer:write_paxos_transfer"
	// ScopeTransferWriteSandboxFiatDeposit: Initiate a test fiat deposit in the Sandbox environment
	ScopeTransferWriteSandboxFiatDeposit Scope = "transfer:write_sandbox_fiat_deposit"
	// ScopeUserAuditLogsReadAuditLog: Retrieve user audit logs
	ScopeUserAuditLogsReadAuditLog Scope = "user_audit_logs:read_audit_log"
	// ScopeWebhooksReadWebhook: Retrieve webhook details and delivery history
	ScopeWebhooksReadWebhook Scope = "webhooks:read_webhook"
)
