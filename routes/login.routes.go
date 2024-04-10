package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	c "github.com/mahl/gotext/config"
	"github.com/mahl/gotext/db"
	m "github.com/mahl/gotext/models"
	ut "github.com/mahl/gotext/utils"
)

func CreateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(c.Config.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := &m.User{}
	json.NewDecoder(r.Body).Decode(reqUser)
	u := &m.User{}

	db.DB.First(u, "name = ?", reqUser.Name)
	if ut.NotFoundCheck(reqUser) {
		message := "Invalid user"
		response := &m.Message{Message: message}
		json.NewEncoder(w).Encode(response)
		return
	}

	if !u.PasswordIsCorrect(reqUser.Password) {
		message := "Invalid user or password"
		response := &m.Message{Message: message}
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := CreateToken(u.ID.String())
	if err != nil {
		response := &m.Message{Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}
