package model

import (
	"time"
)

type Department struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:200;not null"`
	ParentID  *int   `gorm:"foreignKey:ID;check:parent_id <> id"`
	CreatedAt time.Time
	Employees []Employee   `gorm:"foreignKey:DepartmentID;constraint:OnDelete:CASCADE"`
	Children  []Department `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
}

type Employee struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	DepartmentID int
	FullName     string     `gorm:"size:200;not null"`
	Position     string     `gorm:"size:200;not null"`
	HiredAt      *time.Time `gorm:"type:date"`
	CreatedAt    time.Time
}
