package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	NameFull     string    `json:"name_full"`
	Phone        string    `json:"phone"`
	IDNumber     string    `json:"id_number"`
	DateOfBirth  string    `json:"date_of_birth"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AuthToken    string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
