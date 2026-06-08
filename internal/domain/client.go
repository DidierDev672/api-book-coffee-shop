package domain

import "time"

type Client struct {
	ID        string    `json:"id"`
	NameFull  string    `json:"name_full"`
	Phone     string    `json:"phone"`
	Correo    string    `json:"correo"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
