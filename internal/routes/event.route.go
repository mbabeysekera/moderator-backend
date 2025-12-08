package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"coolbreez.lk/moderator/internal/repositories"
	"coolbreez.lk/moderator/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterEventRoutes(routerGroup *gin.RouterGroup, dbPool *pgxpool.Pool) {
	eventRepo := repositories.NewEventRepository(dbPool)
	eventService := services.NewEventService(eventRepo)
	eventsController := controllers.NewEventController(eventService)
	routerGroup.POST("/events/new", eventsController.CreateEvent)
	// routerGroup.GET("/events/:status", eventsController.GetEventsByStatus)
}
