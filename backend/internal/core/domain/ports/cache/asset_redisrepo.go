package cache

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type AssetCacheRepo interface {
	// Set(ctx context.Context, asset *entities.Asset) error
	GetAll(ctx context.Context, page, limit int) (*entities.AssetListCacheData, error)
	SetAll(ctx context.Context, page, limit int, data *entities.AssetListCacheData) error
	// Get(ctx context.Context)
}
