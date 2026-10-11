package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type LocationsRepository interface {
	Create(ctx context.Context, location *entities.Location) (*entities.Location, error)
	Update(ctx context.Context, locationID uuid.UUID, req *entities.UpdateLocation) error
}
