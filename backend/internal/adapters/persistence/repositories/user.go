package repositories

import (
	"context"
	"errors"
	"time"

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

func (r *UserRepository) GetAll(ctx context.Context, page, limit int) ([]*entities.User, int, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Preload("Role").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
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
	if err := r.db.WithContext(ctx).Preload("Role.Permissions").Preload("Role").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&user), nil
}
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Role").First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	return r.modelToEntity(&user), nil
}

func (r *UserRepository) GetEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).Raw(`SELECT EXISTS(SELECT 1) FROM users WHERE email = ?`, email).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, req *entities.UpdateUser) error {
	updates := map[string]interface{}{}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashPassword string) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password", hashPassword).Error
}

func (r *UserRepository) SetRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", token).Error
}
func (r *UserRepository) GetByRefreshToken(ctx context.Context, token string) (*entities.User, error) {
	var userModel models.User
	if err := r.db.WithContext(ctx).First(&userModel, "refresh_token = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token invalid")
		}
		return nil, err
	}
	return r.modelToEntity(&userModel), nil
}

func (r *UserRepository) SetResetToken(ctx context.Context, email, resetToken string) error {
	expiry := time.Now().Add(15 * time.Minute)
	return r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"reset_token":        resetToken,
		"reset_token_expiry": expiry,
	}).Error

}
func (r UserRepository) GetByResetToken(ctx context.Context, token string) (*entities.User, error) {
	var userModel models.User
	if err := r.db.WithContext(ctx).Where("reset_token = ? AND reset_token_expiry > ?", token, time.Now()).Preload("Role").First(&userModel).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&userModel), nil
}

func (r *UserRepository) GetPasswordHash(ctx context.Context, id uuid.UUID) (string, error) {
	var user models.User

	if err := r.db.WithContext(ctx).Select("password").First(&user, "id = ?", id).Error; err != nil {
		return "", err
	}
	return user.Password, nil
}

func (r *UserRepository) ClearResetToken(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"reset_token":        "",
		"reset_token_expiry": nil,
	}).Error
}

func (r *UserRepository) UpdateProfession(ctx context.Context, userID, profesID uuid.UUID) error {
	return r.db.Model(models.User{}).Where("id = ?", userID).Update("profession_id", profesID).Error
}

func (r *UserRepository) AddPermission(ctx context.Context, userID, permissID uuid.UUID) error {
	return r.db.WithContext(ctx).Exec(`INSERT INTO user_permissions(user_id,permission_id) VALUES (?,?)`, userID, permissID).Error
}

func (r *UserRepository) modelToEntity(userModel *models.User) *entities.User {
	user := &entities.User{
		ID:        userModel.ID,
		Email:     userModel.Email,
		Firsname:  userModel.FirstName,
		Lastname:  userModel.LastName,
		Avatar:    userModel.Avatar,
		Active:    userModel.Active,
		RoleID:    userModel.RoleID,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}
	if userModel.RoleID != uuid.Nil && userModel.Role != nil {
		user.Role = &entities.Role{
			ID:          userModel.Role.ID,
			Name:        userModel.Role.Name,
			Description: userModel.Role.Description,
			CreatedAt:   userModel.Role.CreatedAt,
			UpdatedAt:   userModel.Role.UpdatedAt,
		}
	}
	return user
}
