package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseOptions struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func (options DatabaseOptions) GetConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		options.User,
		options.Password,
		options.Host,
		options.Port,
		options.Name,
	)
}
func CreatePool(options DatabaseOptions) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), options.GetConnectionString())
}
