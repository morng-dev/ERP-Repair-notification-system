package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RepairService interface {
	CreateRepair(ctx context.Context, req *entities.RepairRequest) error
	UpdateStatusRepair(ctx context.Context, repairID uuid.UUID, status entities.RepairStatus) error
	GetAllRepair(ctx context.Context, page, limit int) ([]*entities.Repair, *entities.PaginationResponse, error)
}
