package controllers

import (
	"fmt"
	"net/http"
	"time"

	"coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventController struct {
	service *services.EventService
}

func NewEventController(pool *pgxpool.Pool) *EventController {
	return &EventController{
		service: services.NewEventService(pool),
	}
}

func (ev *EventController) CreateEvent(c *gin.Context) {
	var event dto.EventRequest
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusCreated, &dto.ErrorStdResponse{
			Status:  constants.RequestFailed,
			Message: fmt.Sprintf("error from event.controller.create[DATA]: %v", c.Request.Body),
			ErrorID: "ev_0000",
			Details: fmt.Sprintf("ERROR: %v", err),
			Time:    time.Now(),
		})
		return
	}

	eventRes, err := ev.service.CreateEvent(c.Request.Context(), event)
	if err != nil {
		c.JSON(http.StatusCreated, &dto.ErrorStdResponse{
			Status:  constants.RequestFailed,
			Message: "Event creation failed",
			ErrorID: "ev_0001",
			Details: fmt.Sprintf("ERROR: %v", err),
			Time:    time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, eventRes)
}

// func GetEventsByStatus(c *gin.Context) {
// 	status := c.Param("status")
// 	var event models.Event
// 	events, err := event.GetEventsByStatus(c.Request.Context(), constants.EventCode(status))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Something is not right.",
// 			"message": fmt.Sprintf("error from event.controller.getEventsByStatus: %v", c.Request.Body),
// 		})
// 	}
// 	c.JSON(http.StatusOK, events)
// }
