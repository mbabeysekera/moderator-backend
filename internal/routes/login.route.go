package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(routeGroup *gin.RouterGroup, controller *controllers.LoginController) {
	routeGroup.POST("/login", controller.Login)
}
