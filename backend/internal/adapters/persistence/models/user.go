package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Email            string    `gorm:"type:varchar(100);uniqueIndex" json:"email" validate:"required,email"`
	Password         string    `gorm:"type:varchar(100)" json:"-" validate:"required"`
	FirstName        string    `gorm:"type:varchar(100)" json:"first_name" validate:"required"`
	LastName         string    `gorm:"type:varchar(100)" json:"last_name" validate:"required"`
	Avatar           string    `gorm:"type:varchar(255)" json:"avatar"`
	Active           bool      `gorm:"default:true" json:"active"`
	RoleID           uuid.UUID `json:"role_id" validate:"required"`
	Role             Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	RefreshToken     string    `gorm:"type:text" json:"-"`
	ResetToken       string    `gorm:"type:text" json:"-"`
	ResetTokenExpiry time.Time `json:"-"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Assets  []Assets `gorm:"foreignKey:OwnerID" json:"assets,omitempty"`
	Chanals []Chanal `gorm:"many2many:user_chanals;" json:"chanals,omitempty"`
}
