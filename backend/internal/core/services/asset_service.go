package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
)

type AssetService struct {
	assetRepo repositories.AssetsRepositories
}

func (s *AssetService) CreateAsset(ctx context.Context, req *entities.CreateAssetRequest) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	asset := &entities.Asset{
		ID:         id,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		LocationID: req.LocationID,
		OwnerID:    req.OwnerID,
		Image:      req.Image,
	}
	if err := s.assetRepo.Create(ctx, asset); err != nil {
		return err
	}
	return nil
}
