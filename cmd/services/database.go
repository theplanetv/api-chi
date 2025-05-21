package services

import (
	"api-chi/cmd/config"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseService struct {
	*pgxpool.Pool
}

func (s *DatabaseService) Open() error {
	config.LoadDatabaseConfig()
	databaseUrl := config.POSTGRES_URL
	if databaseUrl == "" {
		return errors.New("database url is empty")
	}
	pool, err := pgxpool.New(config.CTX, databaseUrl)
	if err != nil {
		return err
	}

	s.Pool = pool
	return nil
}
