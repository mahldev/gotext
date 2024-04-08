package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	c "github.com/mahl/gotext/config"
	"github.com/mahl/gotext/routes"
)

func userIsAuthenticated(r *http.Request) bool {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return false
	}

	tokenParts := strings.Split(authorizationHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return false
	}
	tokenString := tokenParts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secret := c.Config.SecretKey
		return secret, nil
	})
	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	expTime, _ := claims.GetExpirationTime()
	return expTime.Compare(time.Now()) >= 1
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/login" {
			h.ServeHTTP(w, r)
			return
		}

		if strings.Split(r.URL.Path, "/")[1] == "text" {
			h.ServeHTTP(w, r)
			return
		}

		if !userIsAuthenticated(r) {
			message := "Unauthorized"
			response := &map[string]string{"message": message}
			json.NewEncoder(w).Encode(response)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func InitServer() {
	r := mux.NewRouter()
	r.Use(Middleware)

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users", routes.PostUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", routes.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/{id}", routes.UpdateUserHandler).Methods("UPDATE")
	r.HandleFunc("/users/{id}", routes.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/text", routes.GetAllTextHandler).Methods("GET")
	r.HandleFunc("/text/{n}", routes.GetTextHandler).Methods("GET")
	r.HandleFunc("/login", routes.LoginHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(c.Config.Port, r))
}
