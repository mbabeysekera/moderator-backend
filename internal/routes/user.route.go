package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthUserRoutes(routerGroup *gin.RouterGroup,
	authorizationHandler gin.HandlerFunc, rbacHandler gin.HandlerFunc,
	controller *controllers.UserController) {
	routerGroup.Use(authorizationHandler)
	routerGroup.Use(rbacHandler)
	routerGroup.PATCH("/users/update", controller.UserDetailsUpdate)
}

func RegisterAdminUserRoutes(routerGroup *gin.RouterGroup,
	authorizationHandler gin.HandlerFunc, rbacHandler gin.HandlerFunc,
	controller *controllers.UserController) {
	routerGroup.Use(authorizationHandler)
	routerGroup.Use(rbacHandler)
	routerGroup.GET("/introspect", controller.UserSessionIntrospection)
}
