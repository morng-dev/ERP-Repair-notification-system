package cache

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type UserRedisRepo interface {
	GetAll(ctx context.Context, page, limit int) (*entities.UserListCacheData, error)
	SetAll(ctx context.Context, page, limit int, data *entities.UserListCacheData) error
	DeleteAll(ctx context.Context) error
}
