package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orgstructure/internal/domain/model"
	"orgstructure/internal/logger"
	"strings"
	"time"
)

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
		http.Error(w, "invalid request body: duplicate name of department", http.StatusConflict)
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
