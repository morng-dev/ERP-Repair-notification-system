package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
)

type RoleService struct {
	roleRepo    repositories.RoleRepository
	permissison repositories.PermissionsRepository
}

func (s *RoleService) CreateRole(ctx context.Context, req *entities.Role) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	role := &entities.Role{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}
	exists, err := s.roleRepo.GetByNameExist(ctx, req.Name)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Role already exists")
	}
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return err
	}
	return nil
}
