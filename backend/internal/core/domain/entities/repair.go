package entities

import (
	"time"

	"github.com/google/uuid"
)

type RepairPriority string

const (
	RepairPriorityLow    RepairPriority = "low"
	RepairPriorityNormal RepairPriority = "normal"
	RepairPriorityHigh   RepairPriority = "high"
	RepairPriorityUrgent RepairPriority = "urgent"
)

type RepairStatus string

const (
	RepairStatusPending    RepairStatus = "pending"
	RepairStatusInProgress RepairStatus = "in_progress"
	RepairStatusCompleted  RepairStatus = "completed"
	RepairStatusCancelled  RepairStatus = "cancelled"
)

type Repair struct {
	ID         uuid.UUID      `json:"id"`
	Title      string         `json:"title"`
	AssetID    uuid.UUID      `json:"asset_id"`
	Asset      *Asset         `json:"asset"`
	ReportedBy uuid.UUID      `json:"report_by"`
	Reporter   *User          `json:"reporter"`
	AssignedTo uuid.UUID      `json:"assigned_to"`
	Assignee   *User          `json:"assignee"`
	Status     RepairStatus   `json:"status"`
	Priority   RepairPriority `json:"priority"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type RepairLogs struct {
	ID           uuid.UUID    `json:"id"`
	RequestID    uuid.UUID    `json:"request_id"`
	Request      *Repair      `json:"request"`
	TechnicianID uuid.UUID    `json:"technician_id"`
	Technician   *User        `json:"technician"`
	Description  string       `json:"description"`
	Status       RepairStatus `json:"status"`
	Timestamp    time.Time    `json:"timestamp"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type RepairRequest struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	AssetID    uuid.UUID `json:"asset_id"`
	ReportedBy uuid.UUID `json:"report_by"`
	AssignedTo uuid.UUID `json:"assigned_to"`
}

func (s RepairStatus) IsValid() bool {
	switch s {
	case RepairStatusPending,
		RepairStatusInProgress,
		RepairStatusCompleted,
		RepairStatusCancelled:
		return true
	default:
		return false
	}
}

func CanTransition(current, next RepairStatus) bool {
	switch current {
	case RepairStatusPending:
		return next == RepairStatusInProgress ||
			next == RepairStatusCancelled

	case RepairStatusInProgress:
		return next == RepairStatusCompleted ||
			next == RepairStatusCancelled

	default:
		return false
	}
}
