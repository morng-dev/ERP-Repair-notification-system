package config

import (
	"log"

	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	log.Println("start database seeding...")
	seedRoles(db)
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
