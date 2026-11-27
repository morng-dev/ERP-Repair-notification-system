package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	Description string    `gorm:"type:varchar(100)" json:"description"`
	Image       string    `gorm:"type:varchar(100)" json:"image"`

	ProfessionID *uuid.UUID  `gorm:"type:uuid;index" json:"profession_id"`
	Profession   *Profession `gorm:"foreignKey:ProfessionID;references:ID" json:"profession,omitempty"`

	Assets []Assets `gorm:"foreignKey:CategoryID;references:ID" json:"assets"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
