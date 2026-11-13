package repositories

import (
	"context"

	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repositories.CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *entities.Category) error {
	categoryModel := &models.Category{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
	}
	return r.db.WithContext(ctx).Create(&categoryModel).Error
}

func (r *CategoryRepository) GetAllCategory(ctx context.Context, page, limit int) ([]*entities.Category, int, error) {
	var categoryModel []models.Category
	var total int64
	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&categoryModel).Error; err != nil {
		return nil, 0, err
	}
	var result []*entities.Category
	for _, category := range categoryModel {
		result = append(result, &entities.Category{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}
	return result, int(total), nil
}

func (r *CategoryRepository) GetNameExist(ctx context.Context, name string) (bool, error) {
	var exist bool
	if err := r.db.WithContext(ctx).Raw("SELECT EXISTS(SELECT 1) FROM categories WHERE name = ?", name).Scan(&exist).Error; err != nil {
		return false, err
	}
	return exist, nil
}
