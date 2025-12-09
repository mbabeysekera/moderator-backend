package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"coolbreez.lk/moderator/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(routerGroup *gin.RouterGroup,
	authorizationHandler gin.HandlerFunc, rbacHandler gin.HandlerFunc,
	service *services.UserServiceImpl) {
	userController := controllers.NewUserController(service)
	routerGroup.Use(authorizationHandler)
	routerGroup.Use(rbacHandler)
	routerGroup.PATCH("/user", userController.UserDetailsUpdate)
}
