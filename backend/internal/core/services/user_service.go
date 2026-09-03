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
	userRepo       repositories.UserRepository
	professionRepo repositories.ProfessionRepository
}

func NewUserService(userRepo repositories.UserRepository) services.UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserAll(ctx context.Context, page, limit int) ([]*entities.User, *entities.PaginationResponse, error) {
	users, total, err := s.userRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}
	totalpages := int(math.Ceil(float64(total) / float64(limit)))
	pagination := &entities.PaginationResponse{
		Page:       total,
		Limit:      limit,
		TotalPages: totalpages,
		TotalItems: total,
	}
	return users, pagination, nil
}

func (s *UserService) UserUpdateProfess(ctx context.Context, userID, professID uuid.UUID) error {
	if _, err := s.professionRepo.GetByID(ctx, professID); err != nil {
		return err
	}
	return s.UserUpdateProfess(ctx, userID, professID)

}
