package config

import m "github.com/mahl/gotext/models"

var Config = m.Config{
	Port:      ":8080",
	SecretKey: []byte("my_secret_key"),
}
