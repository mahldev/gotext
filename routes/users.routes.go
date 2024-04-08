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
	ut "github.com/mahl/gotext/utils"
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
		message := "Invalid uuid"
		response := &m.Message{Message: message}
		json.NewEncoder(w).Encode(response)
		return
	}

	db.DB.First(u, id)
	if ut.NotFoundCheck(u) {
		message := fmt.Sprintf("User with id %v not found", id)
		response := &m.Message{Message: message}
		json.NewEncoder(w).Encode(response)
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
		message := "Invalid username"
		response := &m.Message{Message: message}
		json.NewEncoder(w).Encode(response)
		return
	}

	if u.IsValidPassword() {
		w.WriteHeader(http.StatusBadRequest)
		message := "Invalid password"
		response := &m.Message{Message: message}
		json.NewEncoder(w).Encode(response)
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
		response := &m.Message{Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	user := &m.User{}

	db.DB.First(&user, id)
	if ut.NotFoundCheck(user) {
		w.WriteHeader(http.StatusNotFound)
		message := "Error processing request"
		response := &m.Message{Message: message}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	userInRequest := &m.User{}
	json.NewDecoder(r.Body).Decode(userInRequest)

	userInDB := &m.User{}
	db.DB.First(&userInDB, id)
	if ut.NotFoundCheck(userInRequest) {
		w.WriteHeader(http.StatusNotFound)
		message := "Error processing request"
		response := &m.Message{Message: message}
		json.NewEncoder(w).Encode(response)
		return
	}

	userIdUint, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message := "Error processing id request"
		response := &m.Message{Message: message}
		json.NewEncoder(w).Encode(response)
		return
	}

	if userInDB.ID == userIdUint {
		timeNow := time.Now()
		userInDB.Name = userInRequest.Name
		userInDB.Password = userInRequest.Password
		userInDB.UdpateAt = &timeNow
		db.DB.Save(&userInDB)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInDB)
}
