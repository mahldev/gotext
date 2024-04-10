package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Port          string
	DBName        string
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBParams      string
	AuthSecretKey []byte
}

var Config = InitConfig()

func InitConfig() *config {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		log.Fatal("Errr loading .env file")
	}

	return &config{
		Port:          getEnv("SERVER_PORT", "8080"),
		DBName:        getEnv("DB_NAME", "gotext"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", "12345678"),
		DBHost:        getEnv("DB_HOST", "127.0.0.1"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBParams:      getEnv("DB_PARAMS", "parseTime=true"),
		AuthSecretKey: []byte(getEnv("AUTH_SECRET_KEY", "my_secret_key")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
