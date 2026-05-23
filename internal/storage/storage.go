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
	var d model.Department
	result := s.DB.First(&d, id)
	return d, result.Error
}

func (s Storage) GetChildren(id int) []model.Department {
	var children []model.Department
	if id != 0 {
		s.DB.Where("parent_id = ?", id).Find(&children)
	} else {
		s.DB.Where("parent_id IS NULL").Find(&children)
	}
	return children
}

func (s Storage) GetEmployees(id int) []model.Employee {
	var employees []model.Employee
	s.DB.Order("full_name").Where("department_id = ?", id).Find(&employees)
	return employees
}

func (s Storage) UpdateDepartment(d model.Department) (model.Department, error) {
	result := s.DB.Save(&d)
	return d, result.Error
}

func (s Storage) DeleteDepartment(d model.Department) {
	s.DB.Delete(&d)
}

func (s Storage) MoveEmployees(sourceID, destID int) {
	s.DB.Table("employees").Where("department_id = ?", sourceID).Update("department_id", destID)
}
