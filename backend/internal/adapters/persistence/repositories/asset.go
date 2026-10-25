package repositories

import (
	"context"

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

func (r *assetsRepositories) Create(ctx context.Context, asset *entities.Asset) error {
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

	for _, image := range asset.Images {
		assetImages := &models.ImageAssets{
			AssetsID: assetModel.ID,
			ImageURL: image.ImageURL,
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
		return nil, err
	}
	return r.modelToEntity(&asset), nil
}

func (r *assetsRepositories) GetByAssetExist(ctx context.Context, name string) (bool, error) {
	var exists bool

	if err := r.db.WithContext(ctx).Raw(`SELECT EXISTS(SELECT 1) FROM assets WHERE name = ?`, name).Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

func (r *assetsRepositories) GetAll(ctx context.Context, page, limit int) ([]*entities.Asset, int, error) {
	var assets []models.Assets
	var total int64

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Model(&models.Assets{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Preload("Owner").Preload("Category").Preload("Location").Offset(offset).Limit(limit).Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	var result []*entities.Asset
	for _, asset := range assets {
		result = append(result, r.modelToEntity(&asset))
	}
	return result, int(total), nil
}
func (r *assetsRepositories) GetByOwner(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entities.Asset, int, error) {
	var assets []models.Assets
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.Assets{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Where("OwnerID = ?", userID).Find(&assets).Error; err != nil {
		return nil, 0, err
	}
	var result []*entities.Asset

	for _, asset := range assets {
		result = append(result, r.modelToEntity(&asset))
	}

	return result, int(total), nil
}

func (r *assetsRepositories) modelToEntity(assetModel *models.Assets) *entities.Asset {
	asset := &entities.Asset{
		Name:        assetModel.Name,
		Description: assetModel.Description,
		OwnerID:     assetModel.OwnerID,
		Image:       assetModel.Image,
	}
	if assetModel.CategoryID != uuid.Nil && assetModel.Category != nil {
		asset.Category = &entities.Category{
			ID:          assetModel.Category.ID,
			Name:        assetModel.Category.Name,
			Description: assetModel.Category.Description,
			Image:       assetModel.Category.Image,
			CreatedAt:   assetModel.Category.CreatedAt,
			UpdatedAt:   assetModel.Category.UpdatedAt,
		}
	}

	if assetModel.LocationID != uuid.Nil && assetModel.Location != nil {
		asset.Location = &entities.Location{
			ID:          assetModel.Location.ID,
			Description: assetModel.Location.Description,
			Address:     assetModel.Location.Address,
			City:        assetModel.Location.City,
			State:       assetModel.Location.State,
			Latitude:    assetModel.Location.Latitude,
			Longitude:   assetModel.Location.Longitude,
			CreatedAt:   assetModel.Location.CreatedAt,
			UpdatedAt:   assetModel.Location.UpdatedAt,
		}
	}

	return asset
}
