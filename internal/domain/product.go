package domain

import "time"

type Product struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	SupplierID   string    `json:"supplier_id"`
	Name         string    `json:"name"`
	ProductCode  string    `json:"product_code"`
	Categories   []string  `json:"categories"`
	Unit         string    `json:"unit"`
	Quantity     float64   `json:"quantity"`
	MinimumStock float64   `json:"minimum_stock"`
	WineryID     string    `json:"winery_id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
