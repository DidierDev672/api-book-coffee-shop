package domain

import "time"

type Book struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Author          string    `json:"author"`
	Genres          []string  `json:"genres"`
	Photos          []string  `json:"photos"`
	PublicationDate string    `json:"publicationDate"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
