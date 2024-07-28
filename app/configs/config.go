package configs

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type appConfig struct {
	Host    string `env:"APP_HOST"`
	Port    int    `env:"APP_PORT"`
	LogPath string `env:"APP_LOG_PATH"`
}

type databaseConfig struct {
	Host     string `env:"DATABASE_HOST"`
	Port     int    `env:"DATABASE_PORT"`
	Username string `env:"DATABASE_USERNAME"`
	Password string `env:"DATABASE_PASSWORD"`
	Name     string `env:"DATABASE_NAME"`
}

type redisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	Username string `env:"REDIS_USERNAME"`
	Password string `env:"REDIS_PASSWORD"`
}

type config struct {
	App   appConfig
	DB    databaseConfig
	Redis redisConfig
}

var Config config

func InitConfig() error {
	domainFunc := "[] configs.InitConfig ]"
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("%s error loading .env: %v", domainFunc, err)
	}

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("%s error parsing env value: %v", domainFunc, err)
	}

	Config = cfg

	return nil
}