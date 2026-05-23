package web

import (
	"log/slog"
	"orgstructure/internal/domain/model"
)

type Repository interface {
	CreateDepartment(d model.Department) (model.Department, error)
	CreateEmployee(e model.Employee) (model.Employee, error)
	GetDepartment(id int) (model.Department, error)
	GetChildren(id int) []model.Department
	GetEmployees(id int) []model.Employee
	UpdateDepartment(d model.Department) (model.Department, error)
}

type BaseHandler struct {
	data Repository
	log  *slog.Logger
}

func NewBaseHandler(r Repository, l *slog.Logger) BaseHandler {
	return BaseHandler{data: r, log: l}
}
