package domain

import "time"

type Author struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Genres    []string  `json:"genres"`
	BirthDay  string    `json:"birthDay"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
