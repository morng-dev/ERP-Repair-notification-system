package entities

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Image       string    `json:"image"`

	ProfessionID *uuid.UUID  `json:"profession_id"`
	Profession   *Profession `json:"profession,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCategory struct {
	Name         string     `json:"name" validate:"required"`
	Description  string     `json:"description"`
	Image        string     `json:"image"`
	ProfessionID *uuid.UUID `json:"profession_id"`
}
