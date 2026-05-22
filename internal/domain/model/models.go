package model

import "time"

type Department struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:200;not null"`
	ParentID  *int
	CreatedAt time.Time
}

type Employee struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	DepartmentID int
	FullName     string `gorm:"size:200;not null"`
	Position     string `gorm:"size:200;not null"`
	Hired_at     *string
	CreatedAt    time.Time
}
