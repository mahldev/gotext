package api

import (
	"fmt"
	"log"
	"net/http"

	c "github.com/mahl/gotext/config"
)

func InitApi() {
	router := Router()
	port := fmt.Sprintf(":%s", c.Config.Port)

	log.Printf("[INFO] Server listening on port %s.", c.Config.Port)
	log.Fatal(http.ListenAndServe(port, router))
}
