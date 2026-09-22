package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User, password string) (*entities.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetAll(ctx context.Context, page, limit int) ([]*entities.User, int, error)
	Update(ctx context.Context, id uuid.UUID, req *entities.UpdateUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordhash string) error
	SetRefreshToken(ctx context.Context, userID uuid.UUID, token string) error
	GetPasswordHash(ctx context.Context, id uuid.UUID) (string, error)
	UpdateProfession(ctx context.Context, userID, profesID uuid.UUID) error
	AddPermission(ctx context.Context, userID, permissionID uuid.UUID) error
}
