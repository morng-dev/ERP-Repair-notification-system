package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type ProfessionRepository interface {
	Create(ctx context.Context, profession *entities.Profession) (*entities.Profession, error)
	GetAll(ctx context.Context) ([]*entities.Profession, error)
	GetByID(ctx context.Context, professionID uuid.UUID) (*entities.Profession, error)
	GetByName(ctx context.Context, name string) (*entities.Profession, error)
	Update(ctx context.Context, professionID uuid.UUID, req *entities.UpdateProfesion) error
	GetByExists(ctx context.Context, professionID uuid.UUID) (bool, error)
	Delete(ctx context.Context, professionID uuid.UUID) error
	AssignToUser(ctx context.Context, userID, professionID uuid.UUID) error
}
