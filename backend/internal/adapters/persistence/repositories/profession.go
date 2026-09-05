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

func (r *ProfessionRepository) Create(ctx context.Context, req *entities.Profession) (*entities.Profession, error) {
	profess := &models.Profession{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := r.db.WithContext(ctx).Create(&profess).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(profess), nil
}
func (r *ProfessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Profession, error) {
	var profess models.Profession

	if err := r.db.WithContext(ctx).Preload("User").First(&profess, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(&profess), nil
}
func (r *ProfessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Profession{}, "id = ?", id).Error
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

func (r *ProfessionRepository) GetByName(ctx context.Context, name string) (*entities.Profession, error) {
	var profes models.Profession
	if err := r.db.WithContext(ctx).First(&profes, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(&profes), nil
}
func (r *ProfessionRepository) Update(ctx context.Context, id uuid.UUID, req *entities.UpdateProfesion) error {
	updated := map[string]interface{}{}
	if req.Name != "" {
		updated["Name"] = req.Name
	}
	if req.Description != "" {
		updated["Description"] = req.Description
	}
	return r.db.WithContext(ctx).Model(&models.Profession{}).Where("id = ?", id).Updates(updated).Error
}
func (r *ProfessionRepository) modelsToEntities(professModel *models.Profession) *entities.Profession {
	professions := &entities.Profession{
		ID:          professModel.ID,
		Name:        professModel.Name,
		Description: professModel.Description,
		CreatedAt:   professModel.CreatedAt,
		UpdatedAt:   professModel.UpdatedAt,
	}
	for _, profes := range professModel.Users {
		users := entities.User{
			ID:        profes.ID,
			Email:     profes.Email,
			Firsname:  profes.FirstName,
			Lastname:  profes.LastName,
			Avatar:    profes.Avatar,
			Active:    profes.Active,
			CreatedAt: profes.CreatedAt,
			UpdatedAt: profes.UpdatedAt,
		}
		professions.Users = append(professions.Users, users)
	}
	return professions
}
