package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mahl/gotext/routes"
)

var port = ":8080"

func InitServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users", routes.PostUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", routes.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/{id}", routes.UpdateUserHandler).Methods("UPDATE")
	r.HandleFunc("/users/{id}", routes.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/text", routes.GetAllTextHandler).Methods("GET")
	r.HandleFunc("/text/{n}", routes.GetTextHandler).Methods("GET")
	r.HandleFunc("/login/", routes.LoginHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(port, r))
}
