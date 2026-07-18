package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Message   string         `gorm:"type:varchar(255)" json:"message"`
	Type      string         `gorm:"type:varchar(100)" json:"type"`
	Status    string         `gorm:"type:varchar(100)" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
