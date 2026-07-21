package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID   `json:"id"`
	Email        string      `json:"email"`
	Firsname     string      `json:"first_name"`
	Lastname     string      `json:"last_name"`
	Avatar       string      `json:"avatar"`
	Active       bool        `json:"active"`
	ProfessionID uuid.UUID   `json:"profession_id"`
	Profession   *Profession `json:"profession"`
	RoleID       uuid.UUID   `json:"role_id"`
	Role         *Role       `json:"role,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type Profession struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateUser struct {
	FirstName string
	LastName  string
	Avatar    string
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Firsname string `json:"first_name" validate:"required"`
	Lastname string `json:"last_name" validate:"required"`
	Avatar   string `json:"avatar,omitempty"`
}
