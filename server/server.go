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

func Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !userIsAuthenticated(r) {
			message := "Unauthorized"
			response := &map[string]string{"message": message}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	return http.HandlerFunc(fn)
}

func InitServer() {
	r := mux.NewRouter()

	protected := r.PathPrefix("/admin").Subrouter()
	protected.Use(Auth)
	protected.HandleFunc("/users", routes.GetUsersHandler).Methods(http.MethodGet)
	protected.HandleFunc("/users", routes.PostUserHandler).Methods(http.MethodPost)
	protected.HandleFunc("/users/{id}", routes.GetUserHandler).Methods(http.MethodGet)
	protected.HandleFunc("/users/{id}", routes.UpdateUserHandler).Methods(http.MethodPut)
	protected.HandleFunc("/users/{id}", routes.DeleteUserHandler).Methods(http.MethodDelete)

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/text", routes.GetTextHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", routes.LoginHandler).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(c.Config.Port, r))
}
