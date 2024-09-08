package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"env" env-default:"local"`
	BackendDomain string `yaml:"backend_domain" env-default:"local"`
	DatabaseDSN   string `yaml:"database_dsn" env-required:"true"`
	BotToken      string `yaml:"bot_token" env-required:"true"`
}

var (
	instance *Config
)

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("env.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading environment variables: %w", err)
	}
	instance = cfg

	return cfg, nil
}

func GetConfig() *Config {
	return instance
}
