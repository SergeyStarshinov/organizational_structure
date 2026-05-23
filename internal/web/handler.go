package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"orgstructure/internal/domain/model"
	"orgstructure/internal/logger"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
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

func (h BaseHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("start web.GetDepartment")
	w.Header().Set("Content-Type", "application/json")
	idString := strings.TrimSpace(r.PathValue("id"))
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "incorrect department id, must be a number", http.StatusBadRequest)
		h.log.Error("web.GetDepartment, id isn't a number:", logger.Err(err))
		return
	}

	includeEmployees := true
	if strings.TrimSpace(r.URL.Query().Get("include_employees")) == "false" {
		includeEmployees = false
	}
	depth := 1
	depthStr := strings.TrimSpace(r.URL.Query().Get("depth"))
	if depthStr != "" {
		depthQuery, err := strconv.Atoi(depthStr)
		if err != nil || depthQuery < 1 || depthQuery > 5 {
			http.Error(w, "incorrect depth value, must be a number from 1 to 5", http.StatusBadRequest)
			h.log.Error("web.GetDepartment, incorrect depth value: " + depthStr)
			return
		}
		depth = depthQuery
	}
	department, err := h.data.GetDepartment(id)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		h.log.Error("web.GetDepartment, invalid id:", logger.Err(err))
		return
	}
	h.fillDepartmentInfo(&department, depth, includeEmployees)
	h.log.Info(fmt.Sprintf("information about department with ID = %d and include_employees = %t was received",
		department.ID, includeEmployees))
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

	if !h.siblingsCheck(req.Parent_id, req.Name) {
		http.Error(w, "invalid request body: duplicate name of department", http.StatusBadRequest)
		h.log.Error(fmt.Sprintf("web.CreateDepartment, duplicate name %s for parent_id %d:",
			req.Name, req.Parent_id))
		return
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
		newDepartment.ID, newDepartment.Name, parentID))
	departmentInfo, _ := json.Marshal(newDepartment)
	w.Write([]byte(departmentInfo))
}

func (h BaseHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("start web.CreateEmployee")
	w.Header().Set("Content-Type", "application/json")
	idString := strings.TrimSpace(r.PathValue("id"))
	departmentID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "incorrect department id, must be a number", http.StatusBadRequest)
		h.log.Error("web.CreateEmployee, department id isn't a number:", logger.Err(err))
		return
	}
	if _, err := h.data.GetDepartment(departmentID); errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "invalid department id", http.StatusNotFound)
		h.log.Error("web.CreateEmployee, invalid department id: "+idString, logger.Err(err))
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

func (h BaseHandler) ChangeDepartment(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("start web.ChangeDepartment")
	w.Header().Set("Content-Type", "application/json")
	idString := strings.TrimSpace(r.PathValue("id"))
	departmentID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "incorrect department id, must be a number", http.StatusBadRequest)
		h.log.Error("web.ChangeDepartment, department id isn't a number:", logger.Err(err))
		return
	}
	department, err := h.data.GetDepartment(departmentID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "invalid department id", http.StatusNotFound)
		h.log.Error("web.ChangeDepartment, invalid department id: "+idString, logger.Err(err))
		return
	}
	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		h.log.Error("web.ChangeDepartment, invalid request body:", logger.Err(err))
		return
	}

	if reqName, ok := req["name"]; ok {
		if name, ok := reqName.(string); ok {
			department.Name = strings.TrimSpace(name)
		} else {
			http.Error(w, "incorrect name, must be a string", http.StatusBadRequest)
			h.log.Error("web.ChangeDepartment, incorrect name ")
			return
		}
	}
	if reqParentID, ok := req["parent_id"]; ok {
		newParentID := 0
		if parentID, ok := reqParentID.(float64); ok {
			newParentID = int(parentID)
		} else {
			http.Error(w, "incorrect parent_id, must be a number (0 for top level)", http.StatusBadRequest)
			h.log.Error("web.ChangeDepartment, parent_id is not a number")
			return
		}
		if newParentID != 0 {
			if _, err := h.data.GetDepartment(newParentID); errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "invalid parent_id", http.StatusNotFound)
				h.log.Error("web.ChangeDepartment, invalid parent id: "+strconv.Itoa(newParentID), logger.Err(err))
				return
			}
		}
		if !h.siblingsCheck(newParentID, department.Name) {
			http.Error(w, "duplicate name of department for new parent_id", http.StatusBadRequest)
			h.log.Error(fmt.Sprintf("web.ChangeDepartment, duplicate name %s for parent_id %d:",
				department.Name, newParentID))
			return
		}

		if !h.ancestryCheck(department.ID, newParentID) {
			http.Error(w, "incorrect parent_id: cannot move the department inside its subtree",
				http.StatusBadRequest)
			h.log.Error(fmt.Sprintf("web.ChangeDepartment, department %d is an ancestor for %d",
				department.ID, newParentID))
			return
		}
		if newParentID == 0 {
			department.ParentID = nil
		} else {
			department.ParentID = &newParentID
		}
	}

	department, err = h.data.UpdateDepartment(department)
	if err != nil {
		http.Error(w, "updating department error", http.StatusBadRequest)
		h.log.Error("web.ChangeDepartment, updating department error:", logger.Err(err))
		return
	}
	parentID := 0
	if department.ParentID != nil {
		parentID = *department.ParentID
	}
	h.log.Info(fmt.Sprintf("the department has been updated. id %d, parent_id: %d, name: %s",
		department.ID, parentID, department.Name))
	departmentInfo, _ := json.Marshal(department)
	w.Write([]byte(departmentInfo))
}
