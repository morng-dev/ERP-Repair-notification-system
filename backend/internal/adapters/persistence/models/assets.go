package models

import (
	"time"

	"github.com/google/uuid"
)

type Assets struct {
	ID         uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name       string        `gorm:"type:varchar(255)" json:"name" validate:"required"`
	CategoryID uuid.UUID     `json:"category_id" validate:"required"`
	Category   Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	LocationID uuid.UUID     `json:"location_id" validate:"required"`
	Location   Location      `gorm:"foreignKey:LocationID" json:"location"`
	OwnerID    uuid.UUID     `gorm:"type:uuid" json:"owner_id" validate:"required"`
	Owner      User          `gorm:"foreignKey:OwnerID" json:"owner"`
	Image      string        `gorm:"type:varchar(255)" json:"image"`
	Images     []ImageAssets `gorm:"foreignKey:AssetsID" json:"images,omitempty"`
	CreatedAt  time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

type ImageAssets struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AssetsID  uuid.UUID `gorm:"type:uuid" json:"assets_id" validate:"required"`
	ImageURL  string    `gorm:"type:varchar(255)" json:"image_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
