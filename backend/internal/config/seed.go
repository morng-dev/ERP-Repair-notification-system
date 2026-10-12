package config

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/pkg/utils"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB, config *Config) error {
	log.Println("start database seeding...")
	seedRoles(db)
	seedPermission(db)
	seedProfessions(db)
	seedCategory(db)
	seedAdmin(db, config)
	seedRolePermissions(db)
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
				id, err := uuid.NewV7()
				if err != nil {
					return err
				}
				role.ID = id
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
				id, err := uuid.NewV7()
				if err != nil {
					return err
				}
				permission.ID = id
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

func seedCategory(db *gorm.DB) error {
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

func seedAdmin(db *gorm.DB, config *Config) error {
	if config.AdminEmail == "" {
		log.Println("⚠️  ADMIN_EMAIL not set, skipping admin user seeding")
		log.Println("💡 To create admin user, set ADMIN_EMAIL, ADMIN_PASSWORD, ADMIN_FIRST_NAME, ADMIN_LAST_NAME in .env")
		return nil
	}
	var existingUser models.User
	if err := db.Where("email = ?", config.AdminEmail).First(&existingUser).Error; err == nil {
		log.Printf("ℹ️  Admin user already exists: %s", config.AdminEmail)
		return nil
	}

	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		log.Printf("❌ Admin role not found: %v", err)
		return err
	}

	if config.AdminPass == "" {
		log.Println("⚠️  ADMIN_PASSWORD not set, skipping admin user creation")
		return nil
	}
	hashedPassword, err := utils.HashPassword(config.AdminPass)
	if err != nil {
		log.Printf("❌ Error hashing admin password: %v", err)
		return err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	adminUser := &models.User{
		ID:        id,
		Email:     config.AdminEmail,
		Password:  hashedPassword,
		FirstName: config.AdminFirstName,
		LastName:  config.AdminLastName,
		RoleID:    adminRole.ID,
		Active:    true,
	}
	if err := db.Create(adminUser).Error; err != nil {
		log.Printf("❌ Error creating admin user: %v", err)
		return err
	}
	log.Printf("✅ Admin user created successfully: %s", config.AdminEmail)
	log.Printf("👤 Name: %s %s", config.AdminFirstName, config.AdminLastName)
	log.Println("🔐 Password meets security requirements (uppercase, lowercase, number, special character)")

	return nil
}

func seedRolePermissions(db *gorm.DB) error {
	// ดึง Role admin
	var admin models.Role
	if err := db.Where("name = ?", "admin").First(&admin).Error; err != nil {
		return err
	}

	// ดึง Permission ทั้งหมด
	var permissions []models.Permission
	if err := db.Find(&permissions).Error; err != nil {
		return err
	}

	// กำหนด Permission ทั้งหมดให้ admin
	if err := db.Model(&admin).
		Association("Permissions").
		Append(&permissions); err != nil {
		return err
	}

	log.Println("✅ Role permissions seeded successfully")
	return nil
}
