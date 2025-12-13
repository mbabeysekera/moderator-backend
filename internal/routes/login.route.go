package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"coolbreez.lk/moderator/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(routeGroup *gin.RouterGroup, service *services.LoginServiceImpl) {
	loginController := controllers.NewLoginController(service)
	routeGroup.POST("/login", loginController.Login)
}
