package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Mongo struct {
		URI      string
		User     string
		Password string
		Database string
	}
	HTTP struct {
		Host               string
		Port               int
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		MaxHeaderMegabytes int
	}
	GRPC struct {
		Port int
	}
}

func New() (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("mongo", &cfg.Mongo); err != nil {
		return nil, err
	}

	if err := envconfig.Process("grpc", &cfg.GRPC); err != nil {
		return nil, err
	}

	if err := envconfig.Process("http", &cfg.HTTP); err != nil {
		return nil, err
	}

	return cfg, nil
}
