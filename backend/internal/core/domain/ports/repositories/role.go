package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entities.Role) error
	GetAll(ctx context.Context, page, limit int) ([]*entities.Role, int, error)
	GetByName(ctx context.Context, name string) (*entities.Role, error)
	GetByNameExist(ctx context.Context, name string) (bool, error)
	Edit(ctx context.Context, roleID uuid.UUID, req *entities.RoleUpdate) error
	GetByID(ctx context.Context, roleID uuid.UUID) (*entities.Role, error)
	UpdatePermissions(ctx context.Context, roleID, permissionsID uuid.UUID) error
}
