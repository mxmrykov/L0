package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		PG   `yaml:"postgres"`
		Nats `yaml:"nats"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}
	HTTP struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	PG struct {
		Host     string `env-required:"true" yaml:"host" env:"PG_HOST"`
		Port     string `env-required:"true" yaml:"port" env:"PG_PORT"`
		User     string `env-required:"true" yaml:"user" env:"PG_USER"`
		Password string `env-required:"true" yaml:"password" env:"PG_PASSWORD"`
		DBName   string `env-required:"true" yaml:"name" env:"PG_NAME"`
		PgDriver string `env-required:"true" yaml:"pg_driver" env:"PG_PG_DRIVER"`
	}

	Nats struct {
		Host    string `env-required:"true" yaml:"host" env:"NATS_HOST"`
		Port    string `env-required:"true" yaml:"port" env:"NATS_PORT"`
		Cluster string `env-required:"true" yaml:"cluster" env:"NATS_CLUSTER"`
		Client  string `env-required:"true" yaml:"client" env:"NATS_CLIENT"`
		Topic   string `env-required:"true" yaml:"topic" env:"NATS_TOPIC"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	path := "./config/config.yml"

	err := cleanenv.ReadConfig(path, cfg)

	if err != nil {
		return nil, fmt.Errorf("Config error: %v", err)
	}

	err = cleanenv.ReadEnv(cfg)

	if err != nil {
		return nil, fmt.Errorf("Reading enviroment error: %v", err)
	}

	return cfg, nil
}
