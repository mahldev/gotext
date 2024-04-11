package db

import (
	"fmt"
	"log"

	c "github.com/mahl/gotext/config"
	m "github.com/mahl/gotext/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func AutoMigrateTables() {
	err := DB.AutoMigrate(&m.User{})
	if err != nil {
		log.Fatalf("[ERROR] Error migrating User table: %s\n", err.Error())
		return
	}

	log.Println("[INFO] Tables migrations successful.")
}

func InitDBConnection() {
	var err error
	dns := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v",
		c.Config.DBUser,
		c.Config.DBPassword,
		c.Config.DBHost,
		c.Config.DBPort,
		c.Config.DBName,
		c.Config.DBParams,
	)

	DB, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	AutoMigrateTables()

	log.Println("[INFO] DB connection successful.")
}
