package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func LoadFromFile(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
