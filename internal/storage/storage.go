package storage

import (
	"fmt"
	"orgstructure/internal/config"
	"orgstructure/internal/domain/model"

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

	err = db.AutoMigrate(&model.Department{}, &model.Employee{})
	if err != nil {
		return nil, fmt.Errorf("storage.New, migrate: %w", err)
	}

	storage.DB = db
	return &storage, nil
}

func (s Storage) CreateDepartment(d model.Department) (model.Department, error) {
	result := s.DB.Create(&d)
	return d, result.Error
}

func (s Storage) CreateEmployee(e model.Employee) (model.Employee, error) {
	result := s.DB.Create(&e)
	return e, result.Error
}

func (s Storage) GetDepartment(id int) (model.Department, error) {
	// s.DB.Find("id")
	// query := `SELECT id, name
	// 	FROM departments WHERE id=@id`
	// args := pgx.NamedArgs{
	// 	"id": id,
	// }
	// queryRow := r.store.Data.QueryRow(context.Background(), query, args)

	// var department model.Department
	// err := queryRow.Scan(&department.ID, &department.name)
	// if errors.Is(err, pgx.ErrNoRows) {
	// 	return playerDTO, fmt.Errorf("department with id %s not found", id)
	// }
	return model.Department{}, nil
}
