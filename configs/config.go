package configs

import (
	"log"

	"github.com/caarlos0/env/v11"
)

const (
	ModeDev  = "dev"
	ModeProd = "prod"
)

type Config struct {
	Mode      string `env:"MODE"`
	DBURL     string `env:"DB_URL"`
	SecretKey string `env:"SECRET_KEY"`
}

var Cfg Config

func InitConfig() {
	err := env.Parse(&Cfg)
	if err != nil {
		log.Fatal(err)
	}
}
