package domain

import "time"

type Exercise struct {
	ID           string    `json:"id"`
	EquipmentID  string    `json:"equipment_id"`
	Name         string    `json:"name"`
	MuscleGroup  string    `json:"muscle_group"`
	Difficulty   string    `json:"difficulty"`
	VideoURL     string    `json:"video_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
