package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mahl/gotext/auth"
	l "github.com/mahl/gotext/logger"
	"github.com/mahl/gotext/routes"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if r.Method == http.MethodOptions {
			log.Println("[INFO] Handling preflight request")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.Use(l.Logger)
	r.Use(CorsMiddleware)

	protected := r.PathPrefix("/restricted").Subrouter()
	protected.Use(auth.Auth)
	protected.HandleFunc("/users", routes.GetUsersHandler).Methods(http.MethodGet)
	protected.HandleFunc("/users", routes.PostUserHandler).Methods(http.MethodPost)
	protected.HandleFunc("/users/{id}", routes.DeleteUserHandler).Methods(http.MethodDelete)
	protected.HandleFunc("/teststats", routes.SaveTestStatsHandler).Methods(http.MethodPost, http.MethodOptions)
	protected.HandleFunc("/teststats/levelup", routes.ValidatelevelUpTest).Methods(http.MethodPost, http.MethodOptions)
	protected.HandleFunc("/teststats", routes.GetTestStatsHandler).Methods(http.MethodGet, http.MethodOptions)
	protected.HandleFunc("/stats", routes.GetStatsHandler).Methods(http.MethodGet, http.MethodOptions)
	protected.HandleFunc("/profile", routes.GetUserHandler).Methods(http.MethodGet, http.MethodOptions)
	protected.HandleFunc("/profile", routes.UpdateUserHandler).Methods(http.MethodPut, http.MethodOptions)

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/text", routes.GetTextHandler).Methods(http.MethodGet)
	r.HandleFunc("/signup", routes.SignUpHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/login", routes.LoginHandler).Methods(http.MethodPost, http.MethodOptions)

	return r
}
