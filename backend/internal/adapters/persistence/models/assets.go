package models

import (
	"time"

	"github.com/google/uuid"
)

type Assets struct {
	ID         uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string        `gorm:"type:varchar(255);not null" json:"name"`
	CategoryID uuid.UUID     `gorm:"type:uuid;not null" json:"category_id"`
	Category   *Category     `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	LocationID uuid.UUID     `gorm:"type:uuid;not null" json:"location_id"`
	Location   *Location     `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	OwnerID    uuid.UUID     `gorm:"type:uuid;not null" json:"owner_id"`
	Owner      *User         `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Image      string        `gorm:"type:varchar(255)" json:"image,omitempty"`
	Images     []ImageAssets `gorm:"foreignKey:AssetsID" json:"images,omitempty"`
	CreatedAt  time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

type ImageAssets struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	AssetsID  uuid.UUID `gorm:"type:uuid" json:"assets_id" validate:"required"`
	ImageURL  string    `gorm:"type:varchar(255)" json:"image_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
