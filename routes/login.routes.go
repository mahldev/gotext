package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mahl/gotext/auth"
	"github.com/mahl/gotext/db"
	m "github.com/mahl/gotext/models"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	u := &m.User{}
	json.NewDecoder(r.Body).Decode(u)

	if u.IsValidUsername() {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, "Invalid username")
		return
	}

	if u.IsValidPassword() {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, "Invalid password")
		return
	}

	timeNow := time.Now()
	u.ID = uuid.New()
	u.HashPassword()
	u.UpdateAt = &timeNow
	u.CreatedAt = &timeNow
	u.Level = 1

	createdUser := db.DB.Create(u)
	err := createdUser.Error
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		WriteError(w, "User already exits")
		return
	}

	token, err := auth.CreateToken(u.ID.String())
	if err != nil {
		WriteError(w, err.Error())
		return
	}

	response := map[string]string{
		"level": strconv.FormatInt(u.Level, 10),
		"token": token,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
