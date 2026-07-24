package repositories

import (
	"context"

	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repositories.RoleReposittory {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role entities.Role) error {
	roleModel := &models.Role{
		Name:        role.Name,
		Description: role.Description,
	}
	return r.db.WithContext(ctx).Create(roleModel).Error
}

func (r *RoleRepository) modelToEntity(roleModel *models.Role) *entities.Role {
	role := &entities.Role{
		ID:          roleModel.ID,
		Name:        roleModel.Name,
		Description: roleModel.Description,
		CreatedAt:   roleModel.CreatedAt,
		UpdatedAt:   roleModel.UpdatedAt,
	}
	return role
}
