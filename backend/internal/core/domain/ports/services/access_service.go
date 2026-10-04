package services

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type AssetService interface {
	CreateAsset(ctx context.Context, req *entities.CreateAssetRequest) error
}
