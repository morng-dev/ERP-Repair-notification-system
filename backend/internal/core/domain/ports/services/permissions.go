package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type PermissionsService interface {
	CreatePermissions(ctx context.Context, req *entities.Permission) error
	GetAllPermissions(ctx context.Context) ([]*entities.Permission, error)
	UpdatePermissions(ctx context.Context, permissionID uuid.UUID, req *entities.PermissionUpdate) error
	HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
}
