package config

import (
	"log"

	"github.com/caarlos0/env"
)

// PostgresConfig holds the configuration values for PostgreSQL
type Flag struct {
	FlagValue string `env:"FLAG_VALUE" envDefault:"FALSE"`
}

func InitConfig() (*Flag, error) {
	var flagConfig Flag

	// Parse environment variables into the structs
	if err := env.Parse(&flagConfig); err != nil {
		log.Printf("Failed to load Flag config: %v", err)
	}
	return &flagConfig, nil
}
