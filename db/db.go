package db

import (
	"fmt"
	"log"

	m "github.com/mahl/gotext/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func AutoMigrateTables() {
	DB.AutoMigrate(&m.User{})
}

func InitDBConnection() {
	var err error
	dns := "root:12345678@tcp(localhost:3306)/gotext?parseTime=true"

	DB, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	AutoMigrateTables()

	fmt.Println("--------------------------")
	log.Print("DB Connection successfully")
}
