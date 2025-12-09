package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"coolbreez.lk/moderator/config"
	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/db"
	"coolbreez.lk/moderator/internal/logger"
	"coolbreez.lk/moderator/internal/middlewares"
	"coolbreez.lk/moderator/internal/repositories"
	"coolbreez.lk/moderator/internal/routes"
	"coolbreez.lk/moderator/internal/services"
	"coolbreez.lk/moderator/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file not found")
	}

	logger.New()
	slog.Info("Application Starting...")

	ctx := context.Background()
	pool, err := db.InitDB(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("Database initialisation error: %v", err))
	}

	engine := gin.Default()
	// engine.Use(gin.Logger())
	// engine.Use(gin.Recovery())

	basePath, err := config.GetBasePath()
	if err != nil {
		slog.Error(fmt.Sprintf("Application failed to start with error: %v", err))
		panic("critical error")
	}

	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		slog.Error(fmt.Sprintf("Application failed to start with error: %v", err))
		panic("critical error")
	}

	jwtUtil := utils.NewJWTUtil(jwtSecret)
	authorizationHandler := middlewares.AuthorizationHandler
	rbacHandler := middlewares.CheckRBAC

	userRepo := repositories.NewUserRepository(pool)

	signUpService := services.NewSignUpService(userRepo)
	loginService := services.NewLoginService(userRepo, jwtUtil)
	userService := services.NewUserService(userRepo, jwtUtil)

	routerGroup := engine.Group(basePath)
	routes.RegisterHealthCheckRoutes(routerGroup)
	routes.RegisterSignUpRoutes(routerGroup, signUpService)
	routes.RegisterLoginRoutes(routerGroup, loginService)
	// routes.RegisterEventRoutes(routerGroup, pool)
	secureRouterGroup := routerGroup.Group("/users")
	routes.RegisterUserRoutes(secureRouterGroup, authorizationHandler(jwtUtil), rbacHandler(enums.RoleUser), userService)

	appPort, err := config.GetServerPort()
	if err != nil {
		slog.Error(fmt.Sprintf("Application failed to start with error: %v", err))
		panic("critical error")
	}
	err = engine.Run(fmt.Sprintf(":%s", appPort))
	if err != nil {
		log.Fatalf("Application failed to start with error: %v", err)
		panic("critical error")
	}
}
