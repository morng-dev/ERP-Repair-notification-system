package services

import (
	"context"
	"errors"

	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
)

type RoleService struct {
	roleRepo    repositories.RoleRepository
	permissison repositories.PermissionsRepository
}

func (s *RoleService) CreateRole(ctx context.Context, req *entities.Role) error {
	exists, err := s.roleRepo.GetByNameExist(ctx, req.Name)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Role already exists")
	}
	if err := s.roleRepo.Create(ctx, req); err != nil {
		return err
	}
	return nil
}
