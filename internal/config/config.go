package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	ServerPort string `env:"SERVER_PORT"`
}

var AppConfig Config

func LoadConfig() error {
	if err := env.Parse(&AppConfig); err != nil {
		return fmt.Errorf("ошибка парсинга .env: %w", err)
	}
	return nil
}
