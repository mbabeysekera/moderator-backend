package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterSignUpRoutes(routerGroup *gin.RouterGroup, controller *controllers.SignUpController) {
	routerGroup.POST("/signup", controller.CreateUser)
}
