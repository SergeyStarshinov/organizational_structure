package web

import (
	"net/http"
)

type OrgStuctureHandler interface {
	CreateDepartment(w http.ResponseWriter, r *http.Request)
	CreateEmployee(w http.ResponseWriter, r *http.Request)
	GetDepartment(w http.ResponseWriter, r *http.Request)
	ChangeDepartment(w http.ResponseWriter, r *http.Request)
	DeleteDepartment(w http.ResponseWriter, r *http.Request)
}

func CreateMux(handler OrgStuctureHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST 	 /departments/", handler.CreateDepartment)
	mux.HandleFunc("GET 	 /departments/{id}", handler.GetDepartment)
	mux.HandleFunc("POST 	 /departments/{id}/employees/", handler.CreateEmployee)
	mux.HandleFunc("PATCH  /departments/{id}", handler.ChangeDepartment)
	mux.HandleFunc("DELETE /departments/{id}", handler.DeleteDepartment)
	return mux
}
