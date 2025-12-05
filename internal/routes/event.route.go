package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterEventRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/events/new", controllers.CreateEvent)
	routerGroup.GET("/events/:status", controllers.GetEventsByStatus)
}
