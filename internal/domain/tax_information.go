package domain

import "time"

type TaxInformation struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	BusinessID          string    `json:"business_id"`
	TaxRegime           string    `json:"tax_regime"`
	VATResponsible      bool      `json:"vat_responsible"`
	WithholdingTaxpayer bool      `json:"withholding_taxpayer"`
	LargeTaxpayer       bool      `json:"large_taxpayer"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}
