// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// TaxFormsService provides access to the TaxForms API.
//
// All tax forms are associated to a single account_id. An account_id can have
// multiple tax forms associated with it. US 1099-B and 1099-Misc tax forms for
// the previous fiscal year will be available in February of the current
// calendar year. Previous years tax forms are also available by request.
type TaxFormsService struct {
	client *Client
}

// ListTaxFormRevisions list Tax Form Revisions
//
// Requires scope: tax:read_tax_form
func (s *TaxFormsService) ListTaxFormRevisions(ctx context.Context, params *ListTaxFormRevisionsParams, opts ...RequestOption) (*ListTaxFormRevisionsResponse, error) {
	return getWithParamsTyped[ListTaxFormRevisionsResponse](s.client, ctx, "/tax/tax-form-revisions", params, opts...)
}

// ListTaxForms list Tax Forms
//
// Requires scope: tax:read_tax_form
func (s *TaxFormsService) ListTaxForms(ctx context.Context, params *ListTaxFormsParams, opts ...RequestOption) (*ListTaxFormsResponse, error) {
	return getWithParamsTyped[ListTaxFormsResponse](s.client, ctx, "/tax/tax-forms", params, opts...)
}
