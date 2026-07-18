package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Location struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Description string         `json:"description"`
	Address     string         `json:"address"`
	City        string         `json:"city"`
	State       string         `json:"state"`
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
