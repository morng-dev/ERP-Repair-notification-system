package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100)" json:"name"`
	Description string    `gorm:"type:varchar(100)" json:"description"`
	Image       string    `gorm:"type:varchar(100)" json:"image"`
	Assets      []Assets  `gorm:"foreignKey:CategoryID" json:"assets"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
