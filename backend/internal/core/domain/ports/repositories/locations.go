package repositories

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type LocationsRepository interface {
	Create(ctx context.Context, req *entities.Location) (*entities.Location, error)
}
