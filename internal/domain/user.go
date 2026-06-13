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
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (u *User) HasPermission(permission string) bool {
	for _, role := range u.Roles {
		if role == permission || role == "admin" {
			return true
		}
	}
	return false
}
