package main

import (
	"context"
	"log"

	"coolbreez.lk/moderator/internal/db"
	"coolbreez.lk/moderator/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file not found")
	}

	ctx := context.Background()
	pool, err := db.InitDB(ctx)
	if err != nil {
		log.Fatalf("Database initialisation error: %v", err)
	}

	host, err := server.New(pool)
	if err != nil {
		log.Fatalf("Application failed to start with error: %v", err)
	}
	host.AppServer.Use(gin.Logger())
	host.AppServer.Use(gin.Recovery())
	host.Start()
}
