package paxostesting

import (
	"time"

	paxos "github.com/paxosglobal/paxos-sdk-go"
)

// TransferFixture returns a *paxos.Transfer for use in mock responses.
// The id parameter sets the transfer ID.
func TransferFixture(id string) *paxos.Transfer {
	return &paxos.Transfer{
		ID:         id,
		CustomerID: "cust_fixture",
		ProfileID:  "prof_fixture",
		Asset:      "BTC",
		Amount:     paxos.MustDecimal("1.50000000"),
		Fee:        paxos.MustDecimal("0.00010000"),
		Total:      paxos.MustDecimal("1.50010000"),
		Direction:  paxos.TransferDirectionDEBIT,
		Type:       paxos.TransferTypeCRYPTOWITHDRAWAL,
		Status:     paxos.TransferStatusCOMPLETED,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

// ProfileFixture returns a *paxos.Profile for use in mock responses.
func ProfileFixture(id string) *paxos.Profile {
	return &paxos.Profile{
		ID:          id,
		Nickname:    "Test Profile",
		Type:        paxos.ProfileTypeNORMAL,
		Description: paxos.String("Fixture profile for testing"),
	}
}

// IdentityFixture returns a *paxos.Identity for use in mock responses.
func IdentityFixture(id string) *paxos.Identity {
	status := paxos.IdentityStatusAPPROVED
	identityType := paxos.IdentityTypePERSON
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()
	return &paxos.Identity{
		ID:            id,
		Status:        &status,
		SummaryStatus: &status,
		Type:          &identityType,
		PersonDetails: &paxos.PersonDetails{
			FirstName:   paxos.String("Jane"),
			LastName:    paxos.String("Doe"),
			DateOfBirth: paxos.String("1990-01-15"),
		},
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}
}

// OrderFixture returns a *paxos.Order for use in mock responses.
func OrderFixture(id string) *paxos.Order {
	market := paxos.MarketBTCUSD
	side := paxos.OrderSideBUY
	orderType := paxos.OrderTypeMARKET
	status := paxos.OrderStatusFILLED
	baseAmount := paxos.MustDecimal("0.10000000")
	price := paxos.MustDecimal("42000.00")
	amountFilled := paxos.MustDecimal("0.10000000")
	createdAt := time.Now().UTC()
	modifiedAt := time.Now().UTC()
	return &paxos.Order{
		ID:           paxos.String(id),
		ProfileID:    paxos.String("prof_fixture"),
		Market:       &market,
		Side:         &side,
		Type:         &orderType,
		Status:       &status,
		BaseAmount:   &baseAmount,
		Price:        &price,
		AmountFilled: &amountFilled,
		CreatedAt:    &createdAt,
		ModifiedAt:   &modifiedAt,
	}
}

// AccountFixture returns a *paxos.Account for use in mock responses.
func AccountFixture(id string) *paxos.Account {
	accountType := paxos.AccountAccountTypeBROKERAGE
	summaryStatus := paxos.IdentityStatusAPPROVED
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()
	return &paxos.Account{
		ID:            paxos.String(id),
		IdentityID:    paxos.String("id_fixture"),
		Type:          &accountType,
		SummaryStatus: &summaryStatus,
		Description:   paxos.String("Fixture account for testing"),
		CreatedAt:     &createdAt,
		UpdatedAt:     &updatedAt,
	}
}

// DepositAddressFixture returns a *paxos.DepositAddress for use in mock responses.
func DepositAddressFixture(id string) *paxos.DepositAddress {
	return &paxos.DepositAddress{
		ID:                    id,
		ProfileID:             "prof_fixture",
		CustomerID:            "cust_fixture",
		CryptoNetwork:         paxos.CryptoNetworkETHEREUM,
		Address:               "0x742d35Cc6634C0532925a3b844Bc9e7595f2bD00",
		ConversionTargetAsset: paxos.DepositAddressConversionTargetAssetNOCONVERSION,
		IdentityID:            paxos.String("id_fixture"),
		AccountID:             paxos.String("acc_fixture"),
		CreatedAt:             time.Now().UTC(),
	}
}

// ListResponse wraps items and an optional cursor in a paginated list response.
// Use this to build multi-page fixture responses.
//
// Example:
//
//	// First page with a cursor pointing to the next page.
//	page1 := paxostesting.ListResponse(
//	    []interface{}{
//	        paxostesting.TransferFixture("txn_1"),
//	        paxostesting.TransferFixture("txn_2"),
//	    },
//	    "cursor_page2",
//	)
//
//	// Last page with no cursor (empty string signals end of results).
//	page2 := paxostesting.ListResponse(
//	    []interface{}{
//	        paxostesting.TransferFixture("txn_3"),
//	    },
//	    "",
//	)
//
// The items are keyed "items", which most list methods read. Methods that read
// another key (addresses, events, payments, statements, issuer_quote_executions)
// need ListResponseWithKey.
func ListResponse(items []interface{}, nextPageCursor string) map[string]interface{} {
	return ListResponseWithKey("items", items, nextPageCursor)
}

// ListResponseWithKey is ListResponse with the items under the given key, for list
// methods whose response does not use "items", for example
// ListResponseWithKey("events", ...) for Events.ListEvents.
func ListResponseWithKey(key string, items []interface{}, nextPageCursor string) map[string]interface{} {
	resp := map[string]interface{}{
		key: items,
	}
	if nextPageCursor != "" {
		resp["next_page_cursor"] = nextPageCursor
	}
	return resp
}
