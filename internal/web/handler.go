package web

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"orgstructure/internal/domain/model"
	"orgstructure/internal/logger"
	"strconv"
	"strings"
	"time"
)

type Repository interface {
	CreateDepartment(d model.Department) (model.Department, error)
	CreateEmployee(e model.Employee) (model.Employee, error)
	GetDepartment(id int) (model.Department, error)
	GetChildren(id int) []model.Department
}

type BaseHandler struct {
	data Repository
	log  *slog.Logger
}

type departmentRequest struct {
	Name      string
	Parent_id int
}

type employeeRequest struct {
	Full_name string
	Position  string
	Hired_at  string
}

func NewBaseHandler(r Repository, l *slog.Logger) BaseHandler {
	return BaseHandler{data: r, log: l}
}

func (h BaseHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("start web.GetDepartment")
	w.Header().Set("Content-Type", "application/json")
	idString := r.PathValue("id")
	id, err := strconv.Atoi(strings.TrimSpace(idString))
	if err != nil {
		http.Error(w, "incorrect department id, must be a number", http.StatusBadRequest)
		h.log.Error("web.GetDepartment, id isn't a number:", logger.Err(err))
		return
	}
	// TODO: add depth and include_employee check
	department, err := h.data.GetDepartment(id)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		h.log.Error("web.GetDepartment, invalid id:", logger.Err(err))
		return
	}
	department.Children = h.data.GetChildren(department.ID)
	h.log.Info(fmt.Sprintf("information about department with ID = %d was received", department.ID))
	departmentInfo, _ := json.Marshal(department)
	w.Write([]byte(departmentInfo))
}

func (h BaseHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("start web.CreateDepartment")
	w.Header().Set("Content-Type", "application/json")
	var req departmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		h.log.Error("web.CreateDepartment, invalid request body:", logger.Err(err))
		return
	}
	req.Name = strings.TrimSpace(req.Name)

	siblings := h.data.GetChildren(req.Parent_id)
	for _, s := range siblings {
		if s.Name == req.Name {
			http.Error(w, "invalid request body: duplicate name of department", http.StatusBadRequest)
			h.log.Error(fmt.Sprintf("web.CreateDepartment, duplicate name %s for parent_id %d:",
				req.Name, req.Parent_id))
			return
		}
	}

	var newDepartment model.Department
	newDepartment.Name = req.Name
	parentID := req.Parent_id
	if parentID != 0 {
		newDepartment.ParentID = &parentID
	}
	newDepartment.CreatedAt = time.Now()
	newDepartment, err := h.data.CreateDepartment(newDepartment)
	if err != nil {
		http.Error(w, "creating department error", http.StatusBadRequest)
		h.log.Error("web.CreateDepartment, creating department error:", logger.Err(err))
		return
	}
	h.log.Info(fmt.Sprintf("new department created. id: %d, name: %s, parent_id: %d",
		newDepartment.ID, newDepartment.Name, *newDepartment.ParentID))
	departmentInfo, _ := json.Marshal(newDepartment)
	w.Write([]byte(departmentInfo))
}

func (h BaseHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("start web.CreateEmployee")
	w.Header().Set("Content-Type", "application/json")
	idString := r.PathValue("id")
	departmentID, err := strconv.Atoi(strings.TrimSpace(idString))
	if err != nil {
		http.Error(w, "incorrect department id, must be a number", http.StatusBadRequest)
		h.log.Error("web.CreateEmployee, department id isn't a number:", logger.Err(err))
		return
	}
	var req employeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		h.log.Error("web.CreateEmployee, invalid request body:", logger.Err(err))
		return
	}

	var newEmployee model.Employee
	if req.Hired_at != "" {
		hiredDate, err := time.Parse("2006-01-02", req.Hired_at)
		if err != nil {
			http.Error(w, "invalid date format, must be YYYY-MM-DD", http.StatusBadRequest)
			h.log.Error("web.CreateEmployee, invalid date format:", logger.Err(err))
			return
		}
		newEmployee.Hired_at = &hiredDate
	}
	newEmployee.DepartmentID = departmentID
	newEmployee.FullName = strings.TrimSpace(req.Full_name)
	newEmployee.Position = strings.TrimSpace(req.Position)
	newEmployee.CreatedAt = time.Now()
	newEmployee, err = h.data.CreateEmployee(newEmployee)
	if err != nil {
		http.Error(w, "creating employee error", http.StatusBadRequest)
		h.log.Error("web.CreateEmployee, creating employee error:", logger.Err(err))
		return
	}
	h.log.Info(fmt.Sprintf("new employee created. id %d, department id: %d, fullname: %s, position %s, hired at %v",
		newEmployee.ID, newEmployee.DepartmentID, newEmployee.FullName, newEmployee.Position, newEmployee.Hired_at))
	employeeInfo, _ := json.Marshal(newEmployee)
	w.Write([]byte(employeeInfo))
}
