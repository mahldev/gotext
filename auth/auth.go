package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mahl/gotext/config"
)

var SecretKey = config.Config.SecretKey

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
		return SecretKey, nil
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
