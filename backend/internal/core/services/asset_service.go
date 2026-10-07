package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"gorm.io/gorm"
)

type AssetService struct {
	assetRepo    repositories.AssetsRepositories
	locationRepo repositories.LocationsRepository
}

func NewAssetsService(assetRepo repositories.AssetsRepositories) services.AssetService {
	return &AssetService{assetRepo: assetRepo}
}

func (s *AssetService) CreateAsset(ctx context.Context, req *entities.CreateAssetRequest) error {
	assetID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	exist, err := s.assetRepo.GetByAssetExist(ctx, req.Name)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("wraning asset already exist !!!!")
	}

	locationID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	location := &entities.Location{
		ID:        locationID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Address:   req.Address,
		City:      req.City,
		State:     req.State,
	}
	createdLocation, err := s.locationRepo.Create(ctx, location)
	if err != nil {
		return err
	}
	var images []entities.ImageAssets

	for _, imgURL := range req.Images {
		images = append(images, entities.ImageAssets{
			ImageURL: imgURL,
		})
	}

	asset := &entities.Asset{
		ID:         assetID,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		LocationID: createdLocation.ID,
		OwnerID:    req.OwnerID,
		Image:      req.Image,
		Images:     images,
	}
	if err := s.assetRepo.Create(ctx, asset); err != nil {
		return err
	}
	return nil
}

func (s *AssetService) GetAssetByID(ctx context.Context, assetID uuid.UUID) (*entities.Asset, error) {
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found asset")
		}
	}
	return asset, nil
}
