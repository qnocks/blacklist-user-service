package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/qnocks/blacklist-user-service/pkg/logger"
	"sync"
)

const envFilename = ".env"

type (
	Config struct {
		Server ServerConfig
		DB     PostgresConfig
	}

	ServerConfig struct {
		Port string `env:"PORT" env-default:"8080"`
	}

	PostgresConfig struct {
		Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"POSTGRES_PORT" env-default:"5432"`
		Username string `env:"POSTGRES_USER" env-default:"root"`
		Password string `env:"POSTGRES_PASS" env-default:"root"`
		DBName   string `env:"POSTGRES_DB" env-default:"postgres"`
		SSLMode  string `env:"POSTGRES_SSL" env-default:"disable"`
	}
)

var instance *Config

func GetConfig() *Config {
	once := sync.Once{}
	once.Do(func() {
		instance = &Config{}

		if err := godotenv.Load(envFilename); err != nil {
			logger.Errorf("error during loading environment variables: %s\n", err.Error())
			return
		}

		if err := cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Errorf("error during mapping environment variables: %s\n", help)
			return
		}
	})

	return instance
}
