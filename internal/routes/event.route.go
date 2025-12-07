package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterEventRoutes(routerGroup *gin.RouterGroup, dbPool *pgxpool.Pool) {
	eventsController := controllers.NewEventController(dbPool)
	routerGroup.POST("/events/new", eventsController.CreateEvent)
	// routerGroup.GET("/events/:status", eventsController.GetEventsByStatus)
}
