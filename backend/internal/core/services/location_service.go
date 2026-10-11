package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
)

type LocationService struct {
	locationRepo repositories.LocationsRepository
}

func NewLocationsService(locationRepo repositories.LocationsRepository) services.LocationService {
	return &LocationService{locationRepo: locationRepo}
}

func (s *LocationService) CreateLocation(ctx context.Context, req *entities.Location) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	location := &entities.Location{
		ID:          id,
		Description: req.Description,
		Address:     req.Address,
		City:        req.City,
		State:       req.State,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}
	if _, err := s.locationRepo.Create(ctx, location); err != nil {
		return err
	}
	return nil
}

func (s *LocationService) UpdateLocation(ctx context.Context, id uuid.UUID, req *entities.UpdateLocation) error {
	return s.locationRepo.Update(ctx, id, req)
}
