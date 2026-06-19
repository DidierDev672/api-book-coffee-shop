package domain

import "time"

type ProductEntryDetail struct {
	Code                string  `json:"code"`
	Product             string  `json:"product"`
	Unit                string  `json:"unit"`
	Quantity            float64 `json:"quantity"`
	UnitCost            float64 `json:"unit_cost"`
	Subtotal            float64 `json:"subtotal"`
	CommercialPolicy    string  `json:"commercial_policy"`
	ProfitMargin        float64 `json:"profit_margin"`
	FixedMarkup         float64 `json:"fixed_markup"`
	SuggestedSellingPrice float64 `json:"suggested_selling_price"`
}

type Deduction struct {
	Code     string  `json:"code"`
	Quantity float64 `json:"quantity"`
}

type FinancialSummary struct {
	PurchaseSubtotal float64 `json:"purchase_subtotal"`
	VAT              float64 `json:"vat"`
	Discount         float64 `json:"discount"`
	PurchaseTotal    float64 `json:"purchase_total"`
}

type ProductEntry struct {
	ID               string              `json:"id"`
	EntryNumber      string              `json:"entry_number"`
	RegisteredDate   string              `json:"registered_date"`
	MovementType     string              `json:"movement_type"`
	Warehouse        string              `json:"warehouse"`
	ResponsibleParty string              `json:"responsible_party"`
	CompanyID        string              `json:"company_id"`
	Details          []ProductEntryDetail `json:"details"`
	FinancialSummary FinancialSummary    `json:"financial_summary"`
	Observations     string              `json:"observations"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
}
