package services

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type RoleService struct {
	roleRepo    repositories.RoleRepository
	permissison repositories.PermissionsRepository
}

func NewRoleService(roleRepo repositories.RoleRepository, permissison repositories.PermissionsRepository) services.RoleService {
	return &RoleService{
		roleRepo:    roleRepo,
		permissison: permissison,
	}
}

func (s *RoleService) CreateRole(ctx context.Context, req *entities.CreateRoleRequest) error {
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

func (s *RoleService) GetAllRole(ctx context.Context, page, limit int) ([]*entities.Role, *entities.PaginationResponse, error) {

	role, total, err := s.roleRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}
	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	pagination := &entities.PaginationResponse{
		Page:       page,
		Limit:      limit,
		TotalPages: totalPage,
		TotalItems: total,
	}
	return role, pagination, nil
}

func (s *RoleService) UpdatePermission(ctc context.Context, roleID, permissionID uuid.UUID) error {
	return s.roleRepo.UpdatePermissions(ctc, roleID, permissionID)
}
