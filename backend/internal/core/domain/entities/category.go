package entities

import "github.com/google/uuid"

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        uuid.UUID `json:"name"`
	Description string    `json:"Description"`
	Image       string    `json:"image"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}
