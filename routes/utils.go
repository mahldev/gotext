package routes

import (
	"encoding/json"
	"net/http"

	m "github.com/mahl/gotext/models"
)

func NotFoundCheck(u *m.User) bool {
	return u.Name == ""
}

func EnableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func WriteError(w http.ResponseWriter, message string) {
	type Message struct {
		Message string `json:"message"`
	}
	response := &Message{Message: message}
	json.NewEncoder(w).Encode(response)
}
