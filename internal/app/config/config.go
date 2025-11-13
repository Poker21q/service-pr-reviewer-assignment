package config

import (
	"errors"
	"os"
)

type (
	Server struct {
		Port string
	}

	Postgres struct {
		User     string
		Password string
		Host     string
		Port     string
		DB       string
	}

	Config struct {
		Server   Server
		Postgres Postgres
	}
)

func Load() (*Config, error) {
	cfg := &Config{
		Server: Server{
			Port: os.Getenv("SERVER_PORT"),
		},
		Postgres: Postgres{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			DB:       os.Getenv("POSTGRES_DB"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	var errs error

	if c.Server.Port == "" {
		errs = errors.Join(errs, errors.New("missing SERVER_PORT environment variable"))
	}
	if c.Postgres.User == "" {
		errs = errors.Join(errs, errors.New("missing POSTGRES_USER environment variable"))
	}
	if c.Postgres.Password == "" {
		errs = errors.Join(errs, errors.New("missing POSTGRES_PASSWORD environment variable"))
	}
	if c.Postgres.Host == "" {
		errs = errors.Join(errs, errors.New("missing POSTGRES_HOST environment variable"))
	}
	if c.Postgres.Port == "" {
		errs = errors.Join(errs, errors.New("missing POSTGRES_PORT environment variable"))
	}
	if c.Postgres.DB == "" {
		errs = errors.Join(errs, errors.New("missing POSTGRES_DB environment variable"))
	}

	return errs
}
