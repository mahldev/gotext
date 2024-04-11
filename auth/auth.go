package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	c "github.com/mahl/gotext/config"
)

var SecretKey = c.Config.AuthSecretKey

func CreateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !userIsAuthenticated(r) {
			message := "Unauthorized"
			response := &map[string]string{"message": message}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r)
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
