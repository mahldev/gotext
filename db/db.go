package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDBConnection() {
	var err error
	dns := "root:12345678@tcp(localhost:3306)/gotext?parseTime=true"

	DB, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("--------------------------")
	log.Print("DB Connection successfully")
}
