package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"coolbreez.lk/moderator/internal/repositories"
	"coolbreez.lk/moderator/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(routerGroup *gin.RouterGroup, pool *pgxpool.Pool) {
	userRepo := repositories.NewUserRepository(pool)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)
	routerGroup.POST("/users/new", userController.CreateUser)
}
