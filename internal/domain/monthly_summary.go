package domain

import "time"

type MonthlySummary struct {
	ID             string    `json:"id"`
	Product        string    `json:"product"`
	BeginningStock float64   `json:"beginning_stock"`
	IncomingOrders float64   `json:"incoming_orders"`
	OutgoingOrders float64   `json:"outgoing_orders"`
	EndingStock    float64   `json:"ending_stock"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
