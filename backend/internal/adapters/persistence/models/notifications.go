package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(100)" json:"title"`
	ReferenceID uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_notification_reference"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_notification_reference" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user"`
	Message     string         `gorm:"type:varchar(255)" json:"message"`
	Type        string         `gorm:"type:varchar(100);uniqueIndex:idx_notification_reference" json:"type"`
	Status      string         `gorm:"type:varchar(100)" json:"status"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
