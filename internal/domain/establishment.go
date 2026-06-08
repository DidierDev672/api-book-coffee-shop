package domain

import "time"

type Establishment struct {
	ID                   string    `json:"id"`
	EstablishmentName    string    `json:"establishment_name"`
	InventoryManager     string    `json:"inventory_manager"`
	WarehousePointOfSale string    `json:"warehouse_point_of_sale"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}
