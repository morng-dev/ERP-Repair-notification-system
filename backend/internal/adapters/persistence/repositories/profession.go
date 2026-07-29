package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type ProfessionRepository struct {
	db *gorm.DB
}

func NewProfessionRepository(db *gorm.DB) repositories.ProfessionRepository {
	return &ProfessionRepository{db: db}
}

func (r *ProfessionRepository) Create(ctx context.Context, req *entities.Profession) error {
	profess := &models.Profession{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := r.db.WithContext(ctx).Create(&profess).Error; err != nil {
		return err
	}
	req.ID = profess.ID
	req.CreatedAt = profess.CreatedAt
	req.UpdatedAt = profess.UpdatedAt
	return nil
}
func (r *ProfessionRepository) AssignToUser(ctx context.Context, userID, profID uuid.UUID) error {
	var profModel models.Profession
	if err := r.db.WithContext(ctx).First(&profModel, "id = ?", profID).Error; err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Model(&models.Profession{}).Where("id = ?", userID).Update("profession_id", profID).Error; err != nil {
		return err
	}
	return nil
}
