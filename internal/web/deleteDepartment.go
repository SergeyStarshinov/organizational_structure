package web

import (
	"errors"
	"fmt"
	"net/http"
	"orgstructure/internal/domain/model"
	"orgstructure/internal/logger"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func (h BaseHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("start web.DeleteDepartment")
	w.Header().Set("Content-Type", "application/json")
	idString := strings.TrimSpace(r.PathValue("id"))
	sourceID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "incorrect department id, must be a number", http.StatusBadRequest)
		h.log.Error("web.DeleteDepartment, id isn't a number:", logger.Err(err))
		return
	}
	mode := strings.TrimSpace(r.URL.Query().Get("mode"))
	if mode != "cascade" && mode != "reassign" {
		http.Error(w, "incorrect mode, must be cascade or reassign", http.StatusBadRequest)
		h.log.Error("web.DeleteDepartment, incorrect mode:" + mode)
		return
	}
	if mode == "reassign" {
		destIDStr := strings.TrimSpace(r.URL.Query().Get("reassign_to_department_id"))
		if destIDStr == "" {
			http.Error(w, "reassign_to_department_id is missing", http.StatusBadRequest)
			h.log.Error("web.DeleteDepartment, reassign mode: reassign_to_department is missing")
			return
		}
		destID, err := strconv.Atoi(destIDStr)
		if err != nil {
			http.Error(w, "incorrect department id, must be a number", http.StatusBadRequest)
			h.log.Error("web.DeleteDepartment, reassign_to_department_id isn't a number:" + destIDStr)
			return
		}
		if _, err := h.data.GetDepartment(destID); errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "invalid reassign_to_department_id", http.StatusNotFound)
			h.log.Error("web.DeleteDepartment, invalid reassign_to_department_id: "+destIDStr, logger.Err(err))
			return
		}
		if children := h.data.GetChildren(sourceID); len(children) != 0 {
			http.Error(w, "cannot delete: the department has children", http.StatusBadRequest)
			h.log.Error("web.DeleteDepartment, department " + idString + " has children")
			return
		}
		h.data.MoveEmployees(sourceID, destID)
	}

	h.data.DeleteDepartment(model.Department{ID: sourceID})
	h.log.Info(fmt.Sprintf("department with ID = %d deleted, mode = %s",
		sourceID, mode))
	w.WriteHeader(http.StatusNoContent)
}
