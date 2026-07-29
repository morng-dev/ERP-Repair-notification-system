package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type ProfessionRepository interface {
	Create(ctx context.Context, profession *entities.Profession) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Profession, error)
	GetByName(ctx context.Context, name string) (*entities.Profession, error)
	Update(ctx context.Context, id uuid.UUID, req *entities.Profession) error
	Delete(ctx context.Context, id uuid.UUID) error
	AssignToUser(ctx context.Context, userID, professionID uuid.UUID) error
}
