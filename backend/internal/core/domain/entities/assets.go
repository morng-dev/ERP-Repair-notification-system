package entities

import (
	"time"

	"github.com/google/uuid"
)

type Asset struct {
	ID         uuid.UUID     `json:"id"`
	Name       string        `json:"name"`
	CategoryID uuid.UUID     `json:"category_id"`
	Category   *Category     `json:"category,omitempty"`
	LocationID uuid.UUID     `json:"location_id"`
	Location   *Location     `json:"location,omitempty"`
	OwnerID    uuid.UUID     `json:"owner_id" validate:"required,uuid"`
	Owner      *User         `json:"owner,omitempty"`
	Image      string        `json:"image,omitempty"`
	Images     []ImageAssets `json:"images,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type ImageAssets struct {
	ID        uuid.UUID `json:"id"`
	AssetsID  uuid.UUID `json:"assets_id" validate:"required"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAssetRequest struct {
	Name       string    `json:"name" validate:"required"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	OwnerID    uuid.UUID `json:"owner_id" validate:"required,uuid"`
	Image      string    `json:"image,omitempty"`
	Images     []string  `json:"images,omitempty"`

	Description string  `json:"description" validate:"required"`
	Address     string  `json:"address" validate:"required"`
	City        string  `json:"city" validate:"required"`
	State       string  `json:"state" validate:"required"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
}
