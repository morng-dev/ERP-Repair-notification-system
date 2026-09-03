package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type UserService interface {
	GetUserAll(ctx context.Context, page, limit int) ([]*entities.User, *entities.PaginationResponse, error)
	// GetUserByID()
	// UpdateUser()
	// DeleteUser()
	UserUpdateProfess(ctx context.Context, userID, profesID uuid.UUID) error
}
