package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repair struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AssetID uuid.UUID `json:"asset_id"`
	Asset   Assets    `gorm:"foreignKey:AssetID" json:"asset,omitempty"`

	ReportedBy uuid.UUID `json:"reported_by"`
	Reporter   User      `gorm:"foreignKey:ReportedBy;references:ID"`

	AssignedTo uuid.UUID `json:"assigned_to"`
	Assignee   User      `gorm:"foreignKey:AssignedTo;references:ID" json:"assignee,omitempty"`

	Status    string         `json:"status"`
	Priority  string         `json:"priority"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type RepairLogs struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RequestID uuid.UUID `json:"request_id"`
	Request   Repair    `gorm:"foreignKey:RequestID" json:"request,omitempty"`

	TechnicianID uuid.UUID `json:"technician_id"`
	Technician   User      `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`

	Description string    `gorm:"type:varchar(100)" json:"description"`
	Status      string    `gorm:"type:varchar(100)" json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
