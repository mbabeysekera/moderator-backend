package db

import (
	"context"
	"log"
	"time"

	"coolbreez.lk/moderator/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// var DB *pgxpool.Pool

func InitDB(ctx context.Context) (*pgxpool.Pool, error) {
	dbConnectionString, err := config.GetDBConnectionString()
	if err != nil {
		return nil, err
	}

	config, err := pgxpool.ParseConfig(dbConnectionString)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	log.Println("✅ Connected to PostgreSQL")
	return pool, nil
}
