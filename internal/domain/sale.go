package domain

import "time"

type SaleDetail struct {
	Code     string  `json:"code"`
	Product  string  `json:"product"`
	Unit     string  `json:"unit"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
}

type Sale struct {
	ID            string       `json:"sale_id"`
	SaleNumber    string       `json:"sale_number"`
	OrderID       string       `json:"order_id"`
	ClientID      string       `json:"client_id"`
	WarehouseID   string       `json:"warehouse_id"`
	OrderType     string       `json:"order_type"`
	Products      []SaleDetail `json:"products"`
	Subtotal      float64      `json:"subtotal"`
	VAT           float64      `json:"vat"`
	Discount      float64      `json:"discount"`
	Total         float64      `json:"total"`
	PaymentMethod string       `json:"payment_method"`
	Status        string       `json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
	CreatedBy     string       `json:"created_by"`
	CompanyID     string       `json:"company_id"`
}
