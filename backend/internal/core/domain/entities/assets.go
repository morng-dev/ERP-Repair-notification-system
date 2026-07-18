package entities

import (
	"time"

	"github.com/google/uuid"
)

type Assets struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	Category  *Category     `json:"category"`
	Location  *Location     `json:"location"`
	OwnerID   uuid.UUID     `json:"owner_id"`
	Owner     *User         `json:"owner"`
	Image     string        `json:"image"`
	Images    []ImageAssets `json:"images,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type ImageAssets struct {
	ID        uuid.UUID `json:"id"`
	AssetsID  uuid.UUID `json:"assets_id"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
