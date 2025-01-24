package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePool() (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), "postgres://postgres:somepass@database:5432/postgres?sslmode=disable")
}
