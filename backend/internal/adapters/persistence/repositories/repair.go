package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type RepairRepository struct {
	db *gorm.DB
}

func NewRepairRepository(db *gorm.DB) repositories.RepairRepository {
	return &RepairRepository{db: db}
}

func (r *RepairRepository) Create(ctx context.Context, req *entities.Repair) error {
	repairModel := &models.Repair{
		ID:         req.ID,
		AssetID:    req.AssetID,
		ReportedBy: req.ReportedBy,
		AssignedTo: req.AssignedTo,
		Status:     req.Status,
		Priority:   req.Priority,
	}
	if err := r.db.WithContext(ctx).Create(repairModel).Error; err != nil {
		return err
	}
	req.CreatedAt = repairModel.CreatedAt
	req.UpdatedAt = repairModel.UpdatedAt
	return nil
}

func (r *RepairRepository) GetByID(ctx context.Context, repairID uuid.UUID) (*entities.Repair, error) {
	var repair models.Repair

	err := r.db.WithContext(ctx).
		Preload("Asset").
		Preload("Asset.Category").
		Preload("Asset.Location").
		Preload("Asset.Owner", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "email")
		}).
		Preload("Asset.Images", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "assets_id", "image_url")
		}).
		First(&repair, "id = ?", repairID).
		Error
	if err != nil {
		return nil, err
	}
	return r.modelsToEntities(&repair), nil
}

// func (r *RepairRepository) GetByAssetID(ctx context.Context, assetID uuid.UUID) (*entities.Repair, error)

func (r *RepairRepository) modelsToEntities(repairModel *models.Repair) *entities.Repair {
	repair := &entities.Repair{
		ID:         repairModel.ID,
		AssetID:    repairModel.AssetID,
		ReportedBy: repairModel.ReportedBy,
		AssignedTo: repairModel.AssignedTo,
		Status:     repairModel.Status,
		Priority:   repairModel.Priority,
		CreatedAt:  repairModel.CreatedAt,
		UpdatedAt:  repairModel.UpdatedAt,
	}

	if repairModel.Asset != nil {
		assetModel := repairModel.Asset

		asset := &entities.Asset{
			ID:         assetModel.ID,
			Name:       assetModel.Name,
			CategoryID: assetModel.CategoryID,
			LocationID: assetModel.LocationID,
			OwnerID:    assetModel.OwnerID,
			Image:      assetModel.Image,
			CreatedAt:  assetModel.CreatedAt,
			UpdatedAt:  assetModel.UpdatedAt,
		}

		if assetModel.Category != nil {
			asset.Category = &entities.Category{
				ID:   assetModel.Category.ID,
				Name: assetModel.Category.Name,
			}
		}

		if assetModel.Location != nil {
			asset.Location = &entities.Location{
				ID:      assetModel.Location.ID,
				Address: assetModel.Location.Address,
			}
		}

		if assetModel.Owner != nil {
			asset.Owner = &entities.User{
				ID:    assetModel.Owner.ID,
				Email: assetModel.Owner.Email,
			}
		}

		if assetModel.Images != nil {
			asset.Images = make([]entities.ImageAssets, 0, len(assetModel.Images))

			for _, imageModel := range assetModel.Images {
				asset.Images = append(asset.Images, entities.ImageAssets{
					ID:       imageModel.ID,
					AssetsID: imageModel.AssetsID,
					ImageURL: imageModel.ImageURL,
				})
			}
		}

		repair.Asset = asset
	}

	return repair
}
