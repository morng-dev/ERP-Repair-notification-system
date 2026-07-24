package repositories

import (
	"context"

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

func (r *roleRepository) Create(ctx context.Context, role entities.Role) error {
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

func (r *roleRepository) modelToEntity(role *models.Role) *entities.Role {
	return &entities.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
