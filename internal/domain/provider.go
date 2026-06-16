package domain

import "time"

type Provider struct {
	ID                string    `json:"id"`
	Code              string    `json:"code"`
	TypePerson        string    `json:"type_person"`
	DocumentType      string    `json:"document_type"`
	DocumentNumber    string    `json:"document_number"`
	VerificationDigit string    `json:"verification_digit"`
	BusinessName      string    `json:"business_name"`
	BusinessActivity  string    `json:"business_activity"`
	Status            bool      `json:"status"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
