package cache

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type PermissionRedisRepo interface {
	Get(ctx context.Context, userID uuid.UUID) ([]*entities.Permission, error)
	Set(ctx context.Context, userID uuid.UUID, permissions []*entities.Permission) error
	Delete(ctx context.Context, userID uuid.UUID) error
}
