package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type PermissionsService struct {
	permissionsRepo repositories.PermissionsRepository
}

func NewPermissionsService(permissionsRepo repositories.PermissionsRepository) services.PermissionsService {
	return &PermissionsService{permissionsRepo: permissionsRepo}
}

func (s *PermissionsService) CreatePermissions(ctx context.Context, req *entities.Permission) error {
	_, exist := s.permissionsRepo.GetByName(ctx, req.Name)
	if exist == nil {
		return errors.New("permission alredy exist")
	}
	return s.permissionsRepo.Create(ctx, req)
}

func (s *PermissionsService) UpdatePermissions(ctx context.Context, permissionID uuid.UUID, req *entities.PermissionUpdate) error {
	return s.permissionsRepo.Update(ctx, permissionID, req)

}
