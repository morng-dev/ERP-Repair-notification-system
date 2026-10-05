package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type AssetService interface {
	CreateAsset(ctx context.Context, req *entities.CreateAssetRequest) error
	GetAssetByID(ctx context.Context, assetID uuid.UUID) (*entities.Asset, error)
}
