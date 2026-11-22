package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RepairRepository interface {
	Create(ctx context.Context, repair *entities.Repair) error
	GetByID(ctx context.Context, repairID uuid.UUID) (*entities.Repair, error)
	GetAll(ctx context.Context, page, limit int) ([]*entities.Repair, int, error)
	UpdateStatus(ctx context.Context, repairID uuid.UUID, status entities.RepairStatus) error
	// GetByAssetID(ctx context.Context, assetID uuid.UUID) (*entities.Repair, error)
}
