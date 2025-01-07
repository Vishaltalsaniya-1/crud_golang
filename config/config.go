package config

import (
	"log"
	"github.com/caarlos0/env"
)

type Flag struct {
	FlagValue string `env:"FLAG_VALUE" envDefault:"TRUE"`
}

func InitConfig() (*Flag, error) {
	var flagConfig Flag

	if err := env.Parse(&flagConfig); err != nil {
		log.Printf("Failed to load Flag config: %v", err)
	}
	return &flagConfig, nil
}
