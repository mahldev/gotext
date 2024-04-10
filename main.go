package main

import (
	"github.com/mahl/gotext/db"
	"github.com/mahl/gotext/server"
)

func Init() {
	db.InitDBConnection()
	server.InitServer()
}

func main() { Init() }
