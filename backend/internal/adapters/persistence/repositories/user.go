package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User, password string) error {
	userModel := &models.User{
		Email:     user.Email,
		Password:  password,
		FirstName: user.Firsname,
		LastName:  user.Lastname,
		Avatar:    user.Avatar,
		Active:    true,
		RoleID:    user.RoleID,
	}
	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, page, limit int) ([]*entities.User, int, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Preload("Role").Preload("Profession").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	var result []*entities.User
	for _, user := range users {
		result = append(result, r.modelToEntity(&user))
	}
	return result, int(total), nil
}
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Role").Preload("Profession").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&user), nil
}
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Role").Preload("Profession").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&user), nil
}
func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, req *entities.UpdateUser) error {
	updates := map[string]interface{}{}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["first_name"] = req.LastName
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}
func (r *UserRepository) modelToEntity(userModel *models.User) *entities.User {
	user := &entities.User{
		ID:           userModel.ID,
		Email:        userModel.Email,
		Firsname:     userModel.FirstName,
		Lastname:     userModel.LastName,
		Avatar:       userModel.Avatar,
		Active:       userModel.Active,
		RoleID:       userModel.RoleID,
		ProfessionID: userModel.ProfessionID,
		CreatedAt:    userModel.CreatedAt,
		UpdatedAt:    userModel.UpdatedAt,
	}
	if userModel.RoleID != uuid.Nil {
		user.Role = &entities.Role{
			ID:          userModel.RoleID,
			Name:        userModel.Role.Name,
			Description: userModel.Role.Description,
			CreatedAt:   userModel.Role.CreatedAt,
			UpdatedAt:   userModel.Role.UpdatedAt,
		}
	}
	if userModel.ProfessionID != uuid.Nil {
		user.Profession = &entities.Profession{
			ID:          userModel.ProfessionID,
			Name:        userModel.Profession.Name,
			Description: userModel.Profession.Description,
			CreatedAt:   userModel.Profession.CreatedAt,
			UpdatedAt:   userModel.Profession.UpdatedAt,
		}
	}
	return user
}
