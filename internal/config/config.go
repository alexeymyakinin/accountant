package config

import (
	"github.com/caarlos0/env"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var config = Config{}

type Config struct {
	DbUrl                 string        `env:"DB_URL,required"`
	JwtDuration           time.Duration `env:"JWT_DURATION,required"`
	JwtSigningKey         string        `env:"JWT_SIGNING_KEY,required"`
	JwtSigningMethodHS256 *jwt.SigningMethodHMAC
}

func LoadConfig() error {
	err := env.Parse(&config)

	config.JwtSigningMethodHS256 = jwt.SigningMethodHS256
	return err
}

func Get() *Config {
	return &config
}
