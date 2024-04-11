package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mahl/gotext/auth"
	l "github.com/mahl/gotext/logger"
	"github.com/mahl/gotext/routes"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.Use(l.Logger)

	protected := r.PathPrefix("/admin").Subrouter()
	protected.Use(auth.Auth)
	protected.HandleFunc("/users", routes.GetUsersHandler).Methods(http.MethodGet)
	protected.HandleFunc("/users", routes.PostUserHandler).Methods(http.MethodPost)
	protected.HandleFunc("/users/{id}", routes.GetUserHandler).Methods(http.MethodGet)
	protected.HandleFunc("/users/{id}", routes.UpdateUserHandler).Methods(http.MethodPut)
	protected.HandleFunc("/users/{id}", routes.DeleteUserHandler).Methods(http.MethodDelete)

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/text", routes.GetTextHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", routes.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/signup", routes.SignupHandler).Methods(http.MethodPost)

	return r
}
