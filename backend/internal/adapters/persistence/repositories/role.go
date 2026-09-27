package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repositories.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *entities.Role) error {
	roleModel := &models.Role{
		Name:        role.Name,
		Description: role.Description,
	}
	return r.db.WithContext(ctx).Create(roleModel).Error
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*entities.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(&role), nil
}

func (r *roleRepository) GetByNameExist(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).Raw(`SELECT EXISTS(SELECT 1) FROM role WHERE name = ?`, name).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *roleRepository) Edit(ctx context.Context, roleID uuid.UUID, req *entities.RoleUpdate) error {
	updates := map[string]interface{}{}

	if req.Name != "" {
		updates["Name"] = req.Name
	}
	if req.Description != "" {
		updates["Description"] = req.Description
	}
	return r.db.WithContext(ctx).Model(models.Role{}).Where("id = ?", roleID).Updates(updates).Error
}

func (r *roleRepository) GetByID(ctx context.Context, roleID uuid.UUID) (*entities.Role, error) {
	var roleModel models.Role
	if err := r.db.WithContext(ctx).Preload("Permissions").First(&roleModel, "id = ?", roleID).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&roleModel), nil
}

func (r *roleRepository) modelToEntity(role *models.Role) *entities.Role {
	return &entities.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
