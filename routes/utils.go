package routes

import (
	"encoding/json"
	"net/http"

	m "github.com/mahl/gotext/models"
)

type Message struct {
	Message string `json:"message"`
}

func NotFoundCheck(u *m.User) bool {
	return u.Name == ""
}

func WriteError(w http.ResponseWriter, message string) {
	response := &Message{Message: message}
	json.NewEncoder(w).Encode(response)
}
