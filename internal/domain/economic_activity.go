package domain

import "time"

type EconomicActivity struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CompanyID   string    `json:"company_id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
