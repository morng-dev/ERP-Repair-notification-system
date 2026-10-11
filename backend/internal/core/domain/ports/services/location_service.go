package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type LocationService interface {
	CreateLocation(ctx context.Context, location *entities.Location) error
	UpdateLocation(ctx context.Context, id uuid.UUID, req *entities.UpdateLocation) error
}
