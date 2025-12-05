package db

import (
	"context"
	"log"
	"time"

	"coolbreez.lk/moderator/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	dbConnectionString, err := config.GetDBConnectionString()
	if err != nil {
		log.Fatal("database connection string error", err)
	}

	config, err := pgxpool.ParseConfig(dbConnectionString)
	if err != nil {
		log.Fatal("Failed to parse DB config:", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal("DB ping failed:", err)
	}

	DB = db
	log.Println("✅ Connected to PostgreSQL")
}
