package services

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type UserService struct {
	userRepo        repositories.UserRepository
	roleRepo        repositories.RoleRepository
	professionRepo  repositories.ProfessionRepository
	PermissionsRepo repositories.PermissionsRepository
}

func NewUserService(userRepo repositories.UserRepository, PermissionsRepo repositories.PermissionsRepository, professionRepo repositories.ProfessionRepository) services.UserService {
	return &UserService{userRepo: userRepo, PermissionsRepo: PermissionsRepo, professionRepo: professionRepo}
}

func (s *UserService) GetUserAll(ctx context.Context, page, limit int) ([]*entities.User, *entities.PaginationResponse, error) {
	users, total, err := s.userRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}
	totalpages := int(math.Ceil(float64(total) / float64(limit)))
	pagination := &entities.PaginationResponse{
		Page:       page,
		Limit:      limit,
		TotalPages: totalpages,
		TotalItems: total,
	}
	return users, pagination, nil
}

func (s *UserService) UserUpdateProfess(ctx context.Context, userID, professID uuid.UUID) error {
	return s.userRepo.UpdateProfession(ctx, userID, professID)

}

func (s *UserService) UserUpdatePermissions(ctx context.Context, userID, permissID uuid.UUID) error {
	return s.userRepo.AddPermission(ctx, userID, permissID)
}

func (s *UserService) UserAddRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.userRepo.AddRole(ctx, userID, roleID)
}
