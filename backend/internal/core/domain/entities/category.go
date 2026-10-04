package entities

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"Description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
