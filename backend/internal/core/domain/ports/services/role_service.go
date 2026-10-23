package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type RoleService interface {
	CreateRole(ctx context.Context, req *entities.CreateRoleRequest) error
	GetAllRole(ctx context.Context, page, limit int) ([]*entities.Role, *entities.PaginationResponse, error)
	UpdatePermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error
}
