package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type PermissionsRepository struct {
	db *gorm.DB
}

func NewPermissionsRepository(db *gorm.DB) repositories.PermissionsRepository {
	return &PermissionsRepository{db: db}
}

func (r *PermissionsRepository) Create(ctx context.Context, permissions *entities.Permission) error {
	permissionModels := models.Permission{
		ID:          permissions.ID,
		Name:        permissions.Name,
		Description: permissions.Description,
	}
	if err := r.db.WithContext(ctx).Create(&permissionModels).Error; err != nil {
		return err
	}
	return nil
}

func (r *PermissionsRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	var permission models.Permission
	if err := r.db.WithContext(ctx).First(&permission, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&permission), nil
}

func (r *PermissionsRepository) GetAll(ctx context.Context, page, limit int) ([]*entities.Permission, int, error) {
	var permissions []models.Permission
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	var result []*entities.Permission
	for _, permission := range permissions {
		result = append(result, r.modelToEntity(&permission))
	}
	return result, int(total), nil
}

func (r *PermissionsRepository) GetByName(ctx context.Context, name string) (*entities.Permission, error) {
	var permission models.Permission
	if err := r.db.WithContext(ctx).First(&permission, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&permission), nil
}

func (r *PermissionsRepository) Update(ctx context.Context, id uuid.UUID, req *entities.Permission) error {
	updated := map[string]interface{}{}

	if req.Name != "" {
		updated["Name"] = req.Name
	}
	if req.Description != "" {
		updated["Description"] = req.Description
	}
	return r.db.WithContext(ctx).Model(&models.Permission{}).Where("id = ?", id).Updates(updated).Error
}

func (r *PermissionsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Permission{}, "id = ?", id).Error
}

func (r *PermissionsRepository) modelToEntity(permissModel *models.Permission) *entities.Permission {
	permission := &entities.Permission{
		ID:          permissModel.ID,
		Name:        permissModel.Name,
		Description: permissModel.Description,
		CreatedAt:   permissModel.CreatedAt,
		UpdatedAt:   permissModel.UpdatedAt,
	}
	return permission
}
