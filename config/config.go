package config

type config struct {
	Port      string
	SecretKey []byte
}

var Config = &config{
	Port:      ":8080",
	SecretKey: []byte("my_secret_key"),
}
