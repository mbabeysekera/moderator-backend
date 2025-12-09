package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"coolbreez.lk/moderator/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterSignUpRoutes(routerGroup *gin.RouterGroup, service *services.SignUpServiceImpl) {
	signUpController := controllers.NewSignUpController(service)
	routerGroup.POST("/signup", signUpController.CreateUser)
}
