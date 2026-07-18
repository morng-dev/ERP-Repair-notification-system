package entities

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
