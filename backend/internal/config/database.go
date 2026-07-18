package config

import (
	"fmt"
	"log"

	"github.com/morng-dev/erp/internal/adapters/persistence/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Setupdatabase(config *Config) *gorm.DB {
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.DBHost, config.DBName, config.DBPass, config.DBPort, config.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database:%s", err)
	}
	runAutoMigrate(db)

	return db
}

func runAutoMigrate(db *gorm.DB) {
	log.Println("start database migration,,,")

	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Repair{},
		&models.RepairLogs{},
		&models.Permission{},
		&models.Notification{},
		&models.Message{},
		&models.Location{},
		&models.ImageAssets{},
		&models.Chanal{},
		&models.Assets{},
		&models.Category{},
	)
	if err != nil {
		log.Fatal("failed to migration...")
	}
	log.Print("database migration successful..")
}
