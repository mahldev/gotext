package routes

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./resources/pages/home.html")
}
