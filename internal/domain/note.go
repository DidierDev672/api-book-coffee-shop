package domain

import "time"

type Note struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Color     string    `json:"color"`
	IDTopic   string    `json:"id_topic"`
	IDBook    string    `json:"id_book"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
