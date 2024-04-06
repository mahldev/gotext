package main

import (
	"github.com/mahl/gotext/db"
	m "github.com/mahl/gotext/models"
	"github.com/mahl/gotext/server"
)

func Init() {
	db.InitDBConnection()
	db.DB.AutoMigrate(&m.User{})
	server.InitServer()
}

func main() {
	Init()
}
