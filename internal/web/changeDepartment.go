package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"orgstructure/internal/logger"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

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
