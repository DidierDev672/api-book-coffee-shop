package domain

import "time"

type ShipmentDetail struct {
	Code     string  `json:"code"`
	Product  string  `json:"product"`
	Unit     string  `json:"unit"`
	Quantity float64 `json:"quantity"`
	UnitCost float64 `json:"unit_cost"`
	Subtotal float64 `json:"subtotal"`
}

type SourceDocument struct {
	EntryID string `json:"entry_id"`
}

type Recipient struct {
	RecipientType string `json:"recipient_type"`
	RecipientID   string `json:"recipient_id"`
}

type ShipmentFinancialSummary struct {
	Subtotal float64 `json:"subtotal"`
	VAT      float64 `json:"vat"`
	Discount float64 `json:"discount"`
	Total    float64 `json:"total"`
}

type Shipment struct {
	ID               string                 `json:"id"`
	ShipmentNumber   string                 `json:"shipment_number"`
	RecordDate       string                 `json:"record_date"`
	MovementType     string                 `json:"movement_type"`
	Status           string                 `json:"status"`
	CompanyID        string                 `json:"company_id"`
	WarehouseID      string                 `json:"warehouse_id"`
	ResponsibleID    string                 `json:"responsible_id"`
	SourceDocument   SourceDocument         `json:"source_document"`
	Recipient        Recipient              `json:"recipient"`
	Details          []ShipmentDetail       `json:"details"`
	FinancialSummary ShipmentFinancialSummary `json:"financial_summary"`
	Remarks          string                 `json:"remarks"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        time.Time              `json:"updatedAt"`
}
