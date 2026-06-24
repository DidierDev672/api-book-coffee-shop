package gin_auth

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex;not null;size:50" json:"name"`
	Users     []User         `gorm:"foreignKey:RoleID" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null;size:100" json:"name"`
	Email     string         `gorm:"uniqueIndex;not null;size:150" json:"email"`
	Password  string         `gorm:"not null;size:255" json:"-"`
	RoleID    uint           `gorm:"not null;index" json:"role_id"`
	Role      Role           `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"role,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
