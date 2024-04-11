package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	params := mux.Vars(r)
	u := &m.User{}
	id, err := uuid.Parse(params["id"])
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
	u.UdpateAt = &timeNow
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
	params := mux.Vars(r)
	id := params["id"]

	userInRequest := &m.User{}
	json.NewDecoder(r.Body).Decode(userInRequest)

	userIdUUID, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteError(w, "Error processing id request")
		return
	}

	userInDB := &m.User{}
	db.DB.First(&userInDB, userIdUUID)
	if NotFoundCheck(userInDB) {
		fmt.Println(userInDB)
		fmt.Println(userInRequest)
		w.WriteHeader(http.StatusNotFound)
		WriteError(w, "User not found")
		return
	}

	if userInDB.ID == userIdUUID {
		timeNow := time.Now()

		if !userInRequest.IsValidUsername() {
			w.WriteHeader(http.StatusBadRequest)
			WriteError(w, "Invalid username")
			return
		}

		if !userInRequest.IsValidPassword() {
			w.WriteHeader(http.StatusBadRequest)
			WriteError(w, "Invalid password")
			return
		}

		userInDB.UdpateAt = &timeNow
		userInDB.Password = userInRequest.Password
		userInDB.Name = userInRequest.Name
		db.DB.Save(&userInDB)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInDB)
}
