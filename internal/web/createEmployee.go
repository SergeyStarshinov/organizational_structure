package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"orgstructure/internal/domain/model"
	"orgstructure/internal/logger"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

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
		newEmployee.HiredAt = &hiredDate
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
		newEmployee.ID, newEmployee.DepartmentID, newEmployee.FullName, newEmployee.Position, newEmployee.HiredAt))
	employeeInfo, _ := json.Marshal(newEmployee)
	w.Write([]byte(employeeInfo))
}
