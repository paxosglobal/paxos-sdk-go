// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

// Services holds all generated API service instances.
type Services struct {
	APICredentials       *APICredentialsService
	AccountMembers       *AccountMembersService
	Accounts             *AccountsService
	AddressVerification  *AddressVerificationService
	Claims               *ClaimsService
	CrossChainTransfer   *CrossChainTransferService
	CryptoDeposits       *CryptoDepositsService
	CryptoWithdrawals    *CryptoWithdrawalsService
	DepositAddresses     *DepositAddressesService
	Events               *EventsService
	Exchange             *ExchangeService
	Fees                 *FeesService
	FiatTransfers        *FiatTransfersService
	Identity             *IdentityService
	IdentityControls     *IdentityControlsService
	IdentityDocuments    *IdentityDocumentsService
	InstitutionMembers   *InstitutionMembersService
	InternalTransfers    *InternalTransfersService
	IssuerQuotes         *IssuerQuotesService
	Limits               *LimitsService
	MarketData           *MarketDataService
	MonitoringAddresses  *MonitoringAddressesService
	OrchestrationRules   *OrchestrationRulesService
	Orchestrations       *OrchestrationsService
	Orders               *OrdersService
	PaxosTransfers       *PaxosTransfersService
	Payments             *PaymentsService
	PayoutGroups         *PayoutGroupsService
	Pricing              *PricingService
	Profiles             *ProfilesService
	QuoteExecutions      *QuoteExecutionsService
	Quotes               *QuotesService
	RewardAddresses      *RewardAddressesService
	Rewards              *RewardsService
	RewardsStatements    *RewardsStatementsService
	SandboxDeposits      *SandboxDepositsService
	SandboxFiatTransfers *SandboxFiatTransfersService
	SandboxIdentity      *SandboxIdentityService
	Settlement           *SettlementService
	StablecoinConversion *StablecoinConversionService
	Statements           *StatementsService
	TaxDocumentation     *TaxDocumentationService
	TaxForms             *TaxFormsService
	Transfers            *TransfersService
	TravelRule           *TravelRuleService
	UserAuditLogs        *UserAuditLogsService
	Webhooks             *WebhooksService
}

// initServices creates and registers all generated API services.
func (c *Client) initServices() Services {
	return Services{
		APICredentials:       &APICredentialsService{client: c},
		AccountMembers:       &AccountMembersService{client: c},
		Accounts:             &AccountsService{client: c},
		AddressVerification:  &AddressVerificationService{client: c},
		Claims:               &ClaimsService{client: c},
		CrossChainTransfer:   &CrossChainTransferService{client: c},
		CryptoDeposits:       &CryptoDepositsService{client: c},
		CryptoWithdrawals:    &CryptoWithdrawalsService{client: c},
		DepositAddresses:     &DepositAddressesService{client: c},
		Events:               &EventsService{client: c},
		Exchange:             &ExchangeService{client: c},
		Fees:                 &FeesService{client: c},
		FiatTransfers:        &FiatTransfersService{client: c},
		Identity:             &IdentityService{client: c},
		IdentityControls:     &IdentityControlsService{client: c},
		IdentityDocuments:    &IdentityDocumentsService{client: c},
		InstitutionMembers:   &InstitutionMembersService{client: c},
		InternalTransfers:    &InternalTransfersService{client: c},
		IssuerQuotes:         &IssuerQuotesService{client: c},
		Limits:               &LimitsService{client: c},
		MarketData:           &MarketDataService{client: c},
		MonitoringAddresses:  &MonitoringAddressesService{client: c},
		OrchestrationRules:   &OrchestrationRulesService{client: c},
		Orchestrations:       &OrchestrationsService{client: c},
		Orders:               &OrdersService{client: c},
		PaxosTransfers:       &PaxosTransfersService{client: c},
		Payments:             &PaymentsService{client: c},
		PayoutGroups:         &PayoutGroupsService{client: c},
		Pricing:              &PricingService{client: c},
		Profiles:             &ProfilesService{client: c},
		QuoteExecutions:      &QuoteExecutionsService{client: c},
		Quotes:               &QuotesService{client: c},
		RewardAddresses:      &RewardAddressesService{client: c},
		Rewards:              &RewardsService{client: c},
		RewardsStatements:    &RewardsStatementsService{client: c},
		SandboxDeposits:      &SandboxDepositsService{client: c},
		SandboxFiatTransfers: &SandboxFiatTransfersService{client: c},
		SandboxIdentity:      &SandboxIdentityService{client: c},
		Settlement:           &SettlementService{client: c},
		StablecoinConversion: &StablecoinConversionService{client: c},
		Statements:           &StatementsService{client: c},
		TaxDocumentation:     &TaxDocumentationService{client: c},
		TaxForms:             &TaxFormsService{client: c},
		Transfers:            &TransfersService{client: c},
		TravelRule:           &TravelRuleService{client: c},
		UserAuditLogs:        &UserAuditLogsService{client: c},
		Webhooks:             &WebhooksService{client: c},
	}
}
