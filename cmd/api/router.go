package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mahl/gotext/auth"
	"github.com/mahl/gotext/services"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	protected := r.PathPrefix("/admin").Subrouter()
	protected.Use(auth.Auth)
	protected.HandleFunc("/users", services.GetUsersHandler).Methods(http.MethodGet)
	protected.HandleFunc("/users", services.PostUserHandler).Methods(http.MethodPost)
	protected.HandleFunc("/users/{id}", services.GetUserHandler).Methods(http.MethodGet)
	protected.HandleFunc("/users/{id}", services.UpdateUserHandler).Methods(http.MethodPut)
	protected.HandleFunc("/users/{id}", services.DeleteUserHandler).Methods(http.MethodDelete)

	r.HandleFunc("/", services.HomeHandler)
	r.HandleFunc("/text", services.GetTextHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", services.LoginHandler).Methods(http.MethodPost)

	return r
}
