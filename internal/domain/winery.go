package domain

import "time"

type Winery struct {
	ID             string    `json:"id"`
	RegisteredDate string    `json:"registered_date"`
	UserID         string    `json:"user_id"`
	CompanyID      string    `json:"company_id"`
	Area           string    `json:"area"`
	Units          string    `json:"units"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
