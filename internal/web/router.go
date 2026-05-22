package web

import (
	"net/http"
)

type OrgStuctureHandler interface {
	CreateDepartment(w http.ResponseWriter, r *http.Request)
	CreateEmployee(w http.ResponseWriter, r *http.Request)
	// GetDepartment(w http.ResponseWriter, r *http.Request)
}

func CreateMux(handler OrgStuctureHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST 	/departments/", handler.CreateDepartment)
	// mux.HandleFunc("GET 	/department/{id}", handler.GetDepartment)
	mux.HandleFunc("POST 	/departments/{id}/employees/", handler.CreateEmployee)

	return mux
}
