package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RoleRepository interface {
	Create(ctx context.Context, role entities.Role) error
	GetByName(ctx context.Context, name string) (*entities.Role, error)
	Edit(ctx context.Context, roleID uuid.UUID, req *entities.RoleUpdate) error
}
