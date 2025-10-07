package database

import (
	"petmatch/internal/config"
	"petmatch/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Pet{},
		&models.AdoptionRequest{},
	)
}
