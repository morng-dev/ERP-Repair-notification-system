package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User, password string) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUser(ctx context.Context, page, limit int) ([]*entities.User, int, error)
	Update(ctx context.Context, id uuid.UUID, req *entities.UpdateUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetPasswordHash(ctx context.Context, id uuid.UUID) (string, error)
}
