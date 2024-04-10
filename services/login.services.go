package services

import (
	"encoding/json"
	"net/http"

	"github.com/mahl/gotext/auth"
	"github.com/mahl/gotext/db"
	m "github.com/mahl/gotext/models"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := &m.User{}
	json.NewDecoder(r.Body).Decode(reqUser)
	u := &m.User{}

	db.DB.First(u, "name = ?", reqUser.Name)
	if NotFoundCheck(reqUser) {
		WriteError(w, "Invalid user")
		return
	}

	if !u.PasswordIsCorrect(reqUser.Password) {
		WriteError(w, "Invalid user or password")
		return
	}

	token, err := auth.CreateToken(u.ID.String())
	if err != nil {
		WriteError(w, err.Error())
		return
	}

	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}
