package config

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	log.Println("start database seeding...")
	seedRoles(db)
	seedPermission(db)
	seedProfessions(db)
	seedCategory(*db)
	log.Println("seeding database success")
	return nil
}

func seedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{
			Name:        "admin",
			Description: "Administrator with full access",
		},
		{
			Name:        "user",
			Description: "Regular user with limited access",
		},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// สร้าง role ใหม่
				if err := db.Create(&role).Error; err != nil {
					log.Printf("❌ Error creating role %s: %v", role.Name, err)
					return err
				}
				log.Printf("✅ Role created: %s", role.Name)
			} else {
				log.Printf("❌ Error checking role %s: %v", role.Name, err)
				return err
			}
		}
	}

	return nil
}

func seedPermission(db *gorm.DB) error {
	permissions := []models.Permission{
		{
			Name:        "user:view",
			Description: "see all user",
		},
		{
			Name:        "user:create",
			Description: "create users",
		},
		{
			Name:        "user:update",
			Description: "edit users",
		},
		{
			Name:        "user:delete",
			Description: "delete user!!!",
		},
		{
			Name:        "role:view",
			Description: "view all roles",
		},
		{
			Name:        "role:create",
			Description: "create roles",
		},
		{
			Name:        "role:update",
			Description: "update roles user",
		},
		{
			Name:        "role:delete",
			Description: "delete role !!!!",
		},
	}

	for _, permission := range permissions {
		var existingPermission models.Permission
		if err := db.Where("name = ?", permission.Name).First(&existingPermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&permission).Error; err != nil {
					log.Printf("❌ Error creating role %s: %v", permission.Name, err)
					return err
				}
				log.Printf("✅ Role created: %s", permission.Name)
			} else {
				log.Printf("❌ Error checking role %s: %v", permission.Name, err)
				return err
			}
		}
	}
	return nil
}

func seedProfessions(db *gorm.DB) error {

	professions := []models.Profession{
		{
			Name:        "IT subport",
			Description: "maintains and troubleshoots an organization's computer systems, networks",
		},
	}

	for _, profession := range professions {
		var exitsProfession models.Profession
		if err := db.Where("name = ?", profession.Name).First(&exitsProfession).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				id, err := uuid.NewV7()
				if err != nil {
					return err
				}

				profession.ID = id
				if err := db.Create(&profession).Error; err != nil {
					log.Printf("❌ Error creating role %s: %v", profession.Name, err)
					return err
				}
				log.Printf("✅ Role created: %s", profession.Name)
			} else {
				log.Printf("❌ Error checking role %s: %v", profession.Name, err)
				return err
			}
		}
	}
	return nil
}

func seedCategory(db gorm.DB) error {
	category := []entities.Category{
		{
			Name:        "เครื่องใช้ไฟฟ้า",
			Description: "สำนักงาน",
		},
	}

	for _, categorys := range category {
		var existCategory models.Category
		if err := db.Where("name = ?", categorys.Name).First(&existCategory).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				id, err := uuid.NewV7()
				if err != nil {
					return err
				}
				categorys.ID = id
				if err := db.Create(&categorys).Error; err != nil {
					log.Printf("Error creating role %s: %v", categorys.Name, err)
					return err
				}
				log.Printf("create category success %v", categorys.Name)
			} else {
				log.Printf("Error category role %s: %v", categorys.Name, err)
				return err
			}
		}
	}
	return nil
}
