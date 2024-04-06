package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	m "github.com/mahl/gotext/models"
	u "github.com/mahl/gotext/utils"
)

func GetAllTextHandler(w http.ResponseWriter, r *http.Request) {
	words := u.ReadWordFile()
	response := &m.Text{Words: words}

	u.EnableCORS(&w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetTextHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	n, err := strconv.ParseUint(params["n"], 0, 0)
	if err != nil {
		http.Error(w, "Invalid value for 'n'", http.StatusBadRequest)
		return
	}

	words := u.ReadWordFileN(n)
	response := &m.Text{Words: words}

	u.EnableCORS(&w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
