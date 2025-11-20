// Package postgres provides functionality for interacting with a PostgreSQL database.
package postgres

import (
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func New(cfg Config) (*PostgresStorage, error) {
	const op = "storage.postgres.NewStrorage"
	return nil, nil
}
