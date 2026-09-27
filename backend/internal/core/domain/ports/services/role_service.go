package services

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RoleService interface {
	CreateRole(ctx context.Context, req *entities.CreateRoleRequest) error
}
