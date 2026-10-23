package services

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/cache"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo      repositories.UserRepository
	UserRedisRepo cache.UserRedisRepo
}

func NewUserService(userRepo repositories.UserRepository, UserRedisRepo cache.UserRedisRepo) services.UserService {
	return &UserService{
		userRepo:      userRepo,
		UserRedisRepo: UserRedisRepo,
	}
}

func (s *UserService) GetUserAll(ctx context.Context, page, limit int) ([]*entities.User, *entities.PaginationResponse, error) {
	cached, err := s.UserRedisRepo.GetAll(ctx, page, limit)

	if err == nil && cached != nil {
		return cached.User, cached.Pagination, nil
	}

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
	cacheData := &entities.UserListCacheData{
		User:       users,
		Pagination: pagination,
	}

	_ = s.UserRedisRepo.SetAll(ctx, page, limit, cacheData)

	return users, pagination, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ไม่พบผู้ใช้")
		}
		return nil, err
	}
	return user, nil
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
