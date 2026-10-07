package repositories

import (
	"context"

	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type LocationsRepository struct {
	db gorm.DB
}

func NewLocationsRepository(db gorm.DB) repositories.LocationsRepository {
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
