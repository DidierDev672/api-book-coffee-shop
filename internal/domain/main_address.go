package domain

import "time"

type MainAddress struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CompanyID  string    `json:"company_id"`
	Country    string    `json:"country"`
	Department string    `json:"department"`
	Address    string    `json:"address"`
	Postcode   string    `json:"postcode"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
