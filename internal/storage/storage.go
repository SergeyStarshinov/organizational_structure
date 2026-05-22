package storage

import (
	"fmt"
	"orgstructure/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
}

func New(cfg *config.Config) (*Storage, error) {
	var storage Storage
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("storage.New: %w", err)
	}

	// err = db.AutoMigrate(&model.Department{}, &model.Employee{})
	// if err != nil {
	// 	return nil, fmt.Errorf("storage.New, migrate: %w", err)
	// }

	storage.DB = db
	return &storage, nil
}
