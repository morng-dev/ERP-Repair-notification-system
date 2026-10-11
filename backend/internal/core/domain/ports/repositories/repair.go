package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RepairRepository interface {
	Create(ctx context.Context, repair *entities.Repair) error
	GetByID(ctx context.Context, repairID uuid.UUID) (*entities.Repair, error)
	GetByAssetID(ctx context.Context, assetID uuid.UUID) (*entities.Repair, error)
}
