package repositories

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RoleRepository interface {
	Create(ctx context.Context, role entities.Role) error
	GetByName(ctx context.Context, name string) (*entities.Role, error)
}
