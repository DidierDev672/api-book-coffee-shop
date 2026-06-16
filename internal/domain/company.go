package domain

import "time"

type Company struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	NIT               string    `json:"nit"`
	SocialReason      string    `json:"social_reason"`
	BusinessName      string    `json:"business_name"`
	TypePerson        string    `json:"type_person"`
	CompanyType       string    `json:"company_type"`
	Status            string    `json:"status"`
	ConstitutionDate  string    `json:"constitution_date"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Cellphone         string    `json:"cellphone"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
