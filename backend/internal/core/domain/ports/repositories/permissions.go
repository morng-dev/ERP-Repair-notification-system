package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type PermissionsRepository interface {
	Create(ctx context.Context, permissions *entities.Permission) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error)
	GetByName(ctx context.Context, name string) (*entities.Permission, error)
	GetAll(ctx context.Context, page, limit int) ([]*entities.Permission, int, error)
	Update(ctx context.Context, id uuid.UUID, req *entities.Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
}
