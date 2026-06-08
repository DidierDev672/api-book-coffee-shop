package domain

import "time"

type OrderDetail struct {
	Code      string  `json:"code"`
	Product   string  `json:"product"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
	Discount  float64 `json:"discount"`
	Taxes     float64 `json:"taxes"`
	Total     float64 `json:"total"`
}

type Order struct {
	ID            string        `json:"id"`
	OrderNumeric  string        `json:"order_numeric"`
	Date          string        `json:"date"`
	Hour          string        `json:"hour"`
	AttendedBy    string        `json:"attended_by"`
	ClientID      string        `json:"client_id"`
	Details       []OrderDetail `json:"details"`
	PaymentMethod string        `json:"payment_method"`
	Status        string        `json:"status"`
	Observations  string        `json:"observations"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}
