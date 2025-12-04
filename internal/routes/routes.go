package routes

import (
	"coolbreez.lk/moderator/config"
	"coolbreez.lk/moderator/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) error {
	basePath, err := config.GetBasePath()
	if err != nil {
		return err
	}
	api := server.Group(basePath)
	api.GET("/health", controllers.HealthCheck)
	api.POST("/event", controllers.CreateEvent)
	return nil
}
