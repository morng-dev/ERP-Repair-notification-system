package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type assetsRepositories struct {
	db *gorm.DB
}

func NewAssetsRepositories(db *gorm.DB) repositories.AssetsRepositories {
	return &assetsRepositories{db: db}
}

func (r *assetsRepositories) Create(ctx context.Context, asset *entities.CreateAssetRequest) error {
	assetModel := &models.Assets{
		ID:         asset.ID,
		Name:       asset.Name,
		CategoryID: asset.CategoryID,
		LocationID: asset.LocationID,
		OwnerID:    asset.OwnerID,
		Image:      asset.Image,
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(assetModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, imageURL := range asset.Images {
		assetImages := &models.ImageAssets{
			AssetsID: assetModel.ID,
			ImageURL: imageURL,
		}
		if err := tx.Create(assetImages).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (r *assetsRepositories) GetByID(ctx context.Context, assetID uuid.UUID) (*entities.Asset, error) {
	var asset models.Assets
	if err := r.db.WithContext(ctx).Preload("Category").First(asset, "id = ?", assetID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found asset data")
		}
		return nil, err
	}
	return r.modelsToEntities(&asset), nil
}

func (r *assetsRepositories) modelsToEntities(assetModel *models.Assets) *entities.Asset {
	asset := &entities.Asset{
		ID:         assetModel.ID,
		Name:       assetModel.Name,
		Image:      assetModel.Image,
		CategoryID: assetModel.CategoryID,
		LocationID: assetModel.LocationID,
		OwnerID:    assetModel.OwnerID,
		CreatedAt:  assetModel.CreatedAt,
		UpdatedAt:  assetModel.UpdatedAt,
	}
	if assetModel.Category.ID != uuid.Nil && assetModel.Category != nil {
		asset.Category = &entities.Category{
			ID:          assetModel.Category.ID,
			Name:        assetModel.Category.Name,
			Description: assetModel.Category.Description,
			Image:       assetModel.Category.Image,
			CreatedAt:   assetModel.Category.CreatedAt,
			UpdatedAt:   asset.Category.UpdatedAt,
		}
	}
	return asset
}
