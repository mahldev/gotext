package server

import (
	"log"
	"net/http"

	c "github.com/mahl/gotext/config"
)

func InitServer() {
	router, port := Router(), c.Config.Port

	log.Fatal(http.ListenAndServe(port, router))
}
