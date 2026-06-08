package domain

import "time"

type Product struct {
	ID             string    `json:"id"`
	ProductCode    string    `json:"product_code"`
	Categories     []string  `json:"categories"`
	Unit           string    `json:"unit"`
	MinimumStock   float64   `json:"minimum_stock"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
