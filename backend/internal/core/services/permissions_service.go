package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/cache"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type PermissionsService struct {
	permissionsRepo      repositories.PermissionsRepository
	permissionsRedisRepo cache.PermissionRedisRepo
}

func NewPermissionsService(permissionsRepo repositories.PermissionsRepository, permissionsRedisRepo cache.PermissionRedisRepo) services.PermissionsService {
	return &PermissionsService{
		permissionsRepo:      permissionsRepo,
		permissionsRedisRepo: permissionsRedisRepo,
	}
}

func (s *PermissionsService) CreatePermissions(ctx context.Context, req *entities.Permission) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	permissions := &entities.Permission{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}
	_, exist := s.permissionsRepo.GetByName(ctx, permissions.Name)
	if exist == nil {
		return errors.New("permission alredy exist")
	}
	if err := s.permissionsRepo.Create(ctx, permissions); err != nil {
		return errors.New("ไม่สามารถสร้างรายการได้")
	}
	return nil
}

func (s *PermissionsService) GetAllPermissions(ctx context.Context) ([]*entities.Permission, error) {
	permission, err := s.permissionsRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return permission, nil

}
func (s *PermissionsService) HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	permissions, err := s.permissionsRedisRepo.Get(ctx, userID)
	if err != nil {
		permissions, err = s.permissionsRepo.GetByUserID(ctx, userID)
		if err != nil {
			return false, err
		}

		_ = s.permissionsRedisRepo.Set(ctx, userID, permissions)

	} else if permissions == nil {
		permissions, err = s.permissionsRepo.GetByUserID(ctx, userID)
		if err != nil {
			return false, err
		}

		_ = s.permissionsRedisRepo.Set(ctx, userID, permissions)
	}

	for _, p := range permissions {
		if p.Name == permission {
			return true, nil
		}
	}
	return false, nil
}

func (s *PermissionsService) UpdatePermissions(ctx context.Context, permissionID uuid.UUID, req *entities.PermissionUpdate) error {
	return s.permissionsRepo.Update(ctx, permissionID, req)
}
