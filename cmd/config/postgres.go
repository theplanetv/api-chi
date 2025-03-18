package config

import (
	"context"
	"os"
)

var CTX context.Context = context.Background()

var (
	// Postgresql database config
	POSTGRES_USERNAME string
	POSTGRES_PASSWORD string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_DATABASE string
	POSTGRES_URL      string
)

func LoadDatabaseConfig() {
	POSTGRES_URL = os.Getenv("POSTGRES_URL")
}
