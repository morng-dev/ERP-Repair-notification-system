package cache

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type AssetCacheRepo interface {
	Set(ctx context.Context, asset *entities.Asset) error
	Get(ctx context.Context)
}
