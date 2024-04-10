package api

import (
	"fmt"
	"log"
	"net/http"

	c "github.com/mahl/gotext/configs"
)

func InitApi() {
	router := Router()
	port := fmt.Sprintf(":%s", c.Config.Port)

	log.Printf("Server listening on port %s.", c.Config.Port)
	log.Fatal(http.ListenAndServe(port, router))
}
