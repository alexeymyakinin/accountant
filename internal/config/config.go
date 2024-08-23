package config

import (
	"github.com/caarlos0/env"
	"time"
)

var config = Config{}

type Config struct {
	DbUrl             string        `env:"DB_URL,required"`
	AuthJwtDuration   time.Duration `env:"JWT_DURATION,required"`
	AuthJwtSigningKey string        `env:"JWT_SIGNING_KEY,required"`
}

func LoadConfig() error {
	return env.Parse(&config)
}

func Get() *Config {
	return &config
}
