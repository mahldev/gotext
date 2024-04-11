package main

import (
	"github.com/mahl/gotext/cmd/api"
	"github.com/mahl/gotext/db"
)

func Init() {
	db.InitDBConnection()
	api.InitApi()
}

func main() {
	Init()
}
