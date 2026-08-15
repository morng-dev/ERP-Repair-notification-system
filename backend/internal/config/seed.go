package config

import (
	"log"

	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	log.Println("start database seeding...")
	seedRoles(db)
	seedPermission(db)
	seedProfessions(db)
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
			Name:        "view_users",
			Description: "see all user",
		},
		{
			Name:        "create_users",
			Description: "create users",
		},
		{
			Name:        "edit_users",
			Description: "edit users",
		},
		{
			Name:        "delete_users",
			Description: "delete user!!!",
		},
		{
			Name:        "view_roles",
			Description: "view all roles",
		},
		{
			Name:        "create_roles",
			Description: "create roles",
		},
		{
			Name:        "edit_roles",
			Description: "update roles user",
		},
		{
			Name:        "delete_role",
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
