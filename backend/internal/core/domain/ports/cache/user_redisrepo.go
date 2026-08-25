package cache

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type UserRedisRepo interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.User, error)
	Set(ctx context.Context, user *entities.User) error
}
