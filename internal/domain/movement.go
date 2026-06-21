package domain

import "time"

type Movement struct {
	ID             string    `json:"id"`
	Date           string    `json:"date"`
	Code           string    `json:"code"`
	Product        string    `json:"product"`
	Unit           string    `json:"unit"`
	Entrance       float64   `json:"entrance"`
	Output         float64   `json:"output"`
	Balance        float64   `json:"balance"`
	UnitCost       float64   `json:"unit_cost"`
	ValorValue     float64   `json:"valor_value"`
	MovementTypeID string    `json:"movement_type_id"`
	Observations   string    `json:"observations"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
