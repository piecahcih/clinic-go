package db

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect(url string) (*sqlx.DB, error) {
	d, err := sqlx.Connect("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	d.SetMaxOpenConns(25)
	d.SetMaxIdleConns(5)
	d.SetConnMaxLifetime(5 * time.Minute)

	return d, nil
}
