package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Firsname  string    `json:"first_name"`
	Lastname  string    `json:"last_name"`
	Avartar   string    `json:"avartar"`
	Phone     string    `json:"phone"`
	Active    bool      `json:"active"`
	RoleID    uuid.UUID `json:"role_id"`
	Role      *Role     `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
