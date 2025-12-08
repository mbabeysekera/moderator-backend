package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterHealthCheckRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/health", controllers.HealthCheck)
}
