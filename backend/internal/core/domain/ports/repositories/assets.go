package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type AssetsRepositories interface {
	Create(ctx context.Context, assets *entities.CreateAssetRequest) error
	GetByID(ctx context.Context, assetID uuid.UUID) (*entities.Asset, error)
}
