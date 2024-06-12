package auth

import (
	"encoding/json"
	"errors"
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

func CorsMiddelware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func GetClaims(r *http.Request) (*jwt.MapClaims, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	tokenParts := strings.Split(authorizationHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, errors.New("authorization header format must be Bearer {token}")
	}
	tokenString := tokenParts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, errors.New("failed to parse token")
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse token claims")
	}

	return &claims, nil
}

func userIsAuthenticated(r *http.Request) bool {
	claims, err := GetClaims(r)
	if err != nil {
		println(err.Error())
		return false
	}

	expTime, _ := claims.GetExpirationTime()
	return expTime.Compare(time.Now()) >= 1
}
