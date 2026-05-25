package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orgstructure/internal/logger"
	"strconv"
	"strings"
)

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
