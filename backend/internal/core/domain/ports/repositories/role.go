package repositories

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RoleReposittory interface {
	Create(ctx context.Context, role entities.Role) error
}
