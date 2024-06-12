package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mahl/gotext/auth"
	"github.com/mahl/gotext/db"
	m "github.com/mahl/gotext/models"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	type UsersResponse struct {
		Users []m.User `json:"users"`
	}

	usersDB := []m.User{}
	db.DB.Find(&usersDB)
	users := &UsersResponse{Users: usersDB}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	u := &m.User{}

	claims, err := auth.GetClaims(r)
	if err != nil {
		WriteError(w, "Unauthorized")
		return
	}

	userID, ok := (*claims)["userID"].(string)
	if !ok {
		WriteError(w, "Invalid token claims")
		return
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		WriteError(w, "Invalid uuid")
		return
	}

	db.DB.First(u, id)
	if NotFoundCheck(u) {
		WriteError(w, fmt.Sprintf("User with id %v not found", id))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
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

	createdUser := db.DB.Create(u)
	err := createdUser.Error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	user := &m.User{}

	userIdUUID, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, "Error processing id request")
		return
	}

	tx := db.DB.Delete(&user, userIdUUID)
	if tx.RowsAffected > 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	u := &m.User{}

	claims, err := auth.GetClaims(r)
	if err != nil {
		WriteError(w, "Unauthorized")
		return
	}

	userID, ok := (*claims)["userID"].(string)
	if !ok {
		WriteError(w, "Invalid token claims")
		return
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		WriteError(w, "Invalid uuid")
		return
	}

	db.DB.First(u, id)
	if NotFoundCheck(u) {
		WriteError(w, fmt.Sprintf("User with id %v not found", id))
		return
	}

	updatedData := &m.User{}
	if err := json.NewDecoder(r.Body).Decode(updatedData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, "Invalid request body")
		return
	}

	if updatedData.IsValidUsername() {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, "Invalid username")
		return
	}

	if updatedData.Password != "" {
		u.Password = updatedData.Password
		u.HashPassword()
	}

	u.Name = updatedData.Name
	u.Email = updatedData.Email
	timeNow := time.Now()
	u.UpdateAt = &timeNow

	if err := db.DB.Save(u).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteError(w, "Failed to update user")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}
