package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBHost             string `env:"DB_HOST"`
	DBPort             string `env:"DB_PORT"`
	DBUser             string `env:"DB_USER"`
	DBPassword         string `env:"DB_PASSWORD"`
	DBName             string `env:"DB_NAME"`
	GrpcPort           string `env:"GRPC_PORT"`
	RestPort           string `env:"REST_PORT"`
	ClickhouseHost     string `env:"CLICKHOUSE_HOST"`
	ClickhousePort     string `env:"CLICKHOUSE_PORT"`
	ClickhouseDatabase string `env:"CLICKHOUSE_DATABASE"`
	ClickhouseUser     string `env:"CLICKHOUSE_USER"`
	ClickhousePassword string `env:"CLICKHOUSE_PASSWORD"`
	RedisHost          string `env:"REDIS_HOST"`
	RedisPort          string `env:"REDIS_PORT"`
	RedisPassword      string `env:"REDIS_PASSWORD"`
	RedisDB            string `env:"REDIS_DB"`
}

var AppConfig Config

func LoadConfig() error {
	if err := env.Parse(&AppConfig); err != nil {
		return fmt.Errorf("ошибка парсинга .env: %w", err)
	}
	return nil
}
