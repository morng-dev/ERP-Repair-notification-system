package entities

import (
	"time"

	"github.com/google/uuid"
)

type RepairRequest struct {
	ID         uuid.UUID `json:"id"`
	AssetID    uuid.UUID `json:"asset_id"`
	ReportedBy uuid.UUID `json:"report_by"`
	Reporter   *User     `json:"reporter"`
	AssignedTo uuid.UUID `json:"assigned_to"`
	Assignee   *User     `json:"assignee"`
	Status     string    `json:"status"`
	Priority   string    `json:"priority"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RepairLogs struct {
	ID           uuid.UUID      `json:"id"`
	RequestID    uuid.UUID      `json:"request_id"`
	Request      *RepairRequest `json:"request"`
	TechnicianID uuid.UUID      `json:"technician_id"`
	Technician   *User          `json:"technician"`
	Description  string         `json:"description"`
	Status       string         `json:"status"`
	Timestamp    time.Time      `json:"timestamp"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
