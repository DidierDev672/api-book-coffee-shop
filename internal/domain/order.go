package domain

import "time"

type OrderDetail struct {
	Code              string  `json:"code"`
	Product           string  `json:"product"`
	Unit              string  `json:"unit"`
	QuantityRequested float64 `json:"quantity_requested"`
	EstimatedCost     float64 `json:"estimated_cost"`
	Subtotal          float64 `json:"subtotal"`
}

type Order struct {
	ID               string           `json:"id"`
	OrderNumeric     string           `json:"order_numeric"`
	OrderType        string           `json:"order_type"`
	Date             string           `json:"date"`
	CompanyID        string           `json:"company_id"`
	UserID           string           `json:"user_id"`
	RequestedBy      string           `json:"requested_by"`
	Details          []OrderDetail    `json:"details"`
	FinancialSummary FinancialSummary `json:"financial_summary"`
	Status           string           `json:"status"`
	ReasonForOrder   string           `json:"reason_for_order"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
}
