package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type LocationsRepository struct {
	db *gorm.DB
}

func NewLocationsRepository(db *gorm.DB) repositories.LocationsRepository {
	return &LocationsRepository{db: db}
}

func (r *LocationsRepository) Create(ctx context.Context, location *entities.Location) (*entities.Location, error) {
	locationModel := &models.Location{
		ID:          location.ID,
		Description: location.Description,
		Address:     location.Address,
		City:        location.City,
		State:       location.State,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
	}

	if err := r.db.WithContext(ctx).Create(locationModel).Error; err != nil {
		return nil, err
	}
	return &entities.Location{
		ID:          locationModel.ID,
		Description: locationModel.Description,
		Address:     locationModel.Address,
		City:        locationModel.City,
		State:       locationModel.State,
		Latitude:    locationModel.Latitude,
		Longitude:   locationModel.Longitude,
		CreatedAt:   locationModel.CreatedAt,
		UpdatedAt:   locationModel.UpdatedAt,
	}, nil
}

func (r *LocationsRepository) Update(ctx context.Context, locationID uuid.UUID, req *entities.UpdateLocation) error {
	var update map[string]interface{}
	if req.Description != "" {
		update["Description"] = req.Description
	}
	if req.City != "" {
		update["City"] = req.City
	}
	if req.Address != "" {
		update["Address"] = req.Address
	}
	if req.State != "" {
		update["State"] = req.State
	}
	return r.db.WithContext(ctx).Model(&models.Location{}).Where("id = ?", locationID).Updates(&update).Error
}
