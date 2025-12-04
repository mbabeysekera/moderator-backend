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
	createEvenTable()
}

func createEvenTable() {
	eventTable := `
		CREATE TABLE IF NOT EXISTS events (
			id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			target TEXT NOT NULL,
			command TEXT NOT NULL,
			status INTEGER NOT NULL,
			createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			processedAt TIMESTAMPTZ NOT NULL 
		)
	`
	_, err := DB.Exec(context.Background(), eventTable)
	if err != nil {
		log.Fatal("Error:", err)
	}
	log.Println("✅ Event Table Created")
}
