package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		Http  HTTP
		Mongo Mongo
		App   App
		Jwt   JWT
	}
	App struct {
		Environment string `envconfig:"env" default:"local"`
	}

	Mongo struct {
		Host               string        `envconfig:"MONGO_HOST" required="true"`
		Port               string        `envconfig:"MONGO_PORT" required="true"`
		User               string        `envconfig:"MONGO_USER" default=""`
		Password           string        `envconfig:"MONGO_PASSWORD" default=""`
		DbName             string        `envconfig:"MONGO_DBNAME" required="true"`
		ConnectionTimeout  time.Duration `envconfig:"MONGO_CONNECTIONTIMEOUT" default="1s"`
		ConnectionAttempts int           `envconfig:"MONGO_CONNECTIONATTEMPTS" default="10"`
		AuthDb             string        `envconfig:"MONGO_AUTHDB" default=""`
	}

	HTTP struct {
		Host            string        `envconfig:"HTTP_HOST" default="localhost"`
		Port            string        `envconfig:"HTTP_PORT" default="80"`
		ReadTimeout     time.Duration `envconfig:"HTTP_READTIMEOUT" default="5s"`
		WriteTimeout    time.Duration `envconfig:"HTTP_WRITETIMEOUT" default="5s"`
		ShutdownTimeout time.Duration `envconfig:"HTTP_SHUTDOWNTIMEOUT" default="3s"`
	}

	JWT struct {
		Secret   string        `envconfig:"JWT_SECRET" required="true"`
		LifeTime time.Duration `envconfig:"JWT_LIFETIME" default="60m"`
	}
)

func GetConfig() (*Config, error) {
	cfg := Config{}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("error with loading env variables: %w", err)
	}

	return &cfg, nil
}
