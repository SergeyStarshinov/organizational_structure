package web

import "orgstructure/internal/domain/model"

type departmentRequest struct {
	Name      string
	Parent_id int
}

type employeeRequest struct {
	Full_name string
	Position  string
	Hired_at  string
}

func (h BaseHandler) fillDepartmentInfo(department *model.Department, depth int, includeEmployees bool) {
	department.Children = h.data.GetChildren(department.ID)
	if includeEmployees {
		department.Employees = h.data.GetEmployees(department.ID)
	}
	if depth > 1 {
		for i := range department.Children {
			h.fillDepartmentInfo(&department.Children[i], depth-1, includeEmployees)
		}
	}
}

func (h BaseHandler) ancestryCheck(departmentID, destID int) bool {
	if destID == 0 {
		return true
	}
	department, _ := h.data.GetDepartment(destID)
	for department.ParentID != nil {
		if *department.ParentID == departmentID {
			return false
		}
		department, _ = h.data.GetDepartment(*department.ParentID)
	}
	return true
}

func (h BaseHandler) siblingsCheck(parentID int, name string) bool {
	siblings := h.data.GetChildren(parentID)
	for _, s := range siblings {
		if s.Name == name {
			return false
		}
	}
	return true
}
