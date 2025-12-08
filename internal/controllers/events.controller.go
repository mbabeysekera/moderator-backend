package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	"github.com/gin-gonic/gin"
)

type EventService interface {
	CreateEvent(rc context.Context, event dto.EventRequest) error
}

type EventController struct {
	service EventService
}

func NewEventController(eventService EventService) *EventController {
	return &EventController{
		service: eventService,
	}
}

func (ev *EventController) CreateEvent(c *gin.Context) {
	var event dto.EventRequest
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorStdResponse{
			Status:  enums.RequestFailed,
			Message: fmt.Sprintf("error from event.controller.create[DATA]: %v", c.Request.Body),
			ErrorID: "ev_0000",
			Details: fmt.Sprintf("ERROR: %v", err),
			Time:    time.Now(),
		})
		return
	}

	err = ev.service.CreateEvent(c.Request.Context(), event)
	if err != nil {
		c.JSON(http.StatusCreated, &dto.ErrorStdResponse{
			Status:  enums.RequestFailed,
			Message: "Event creation failed",
			ErrorID: "ev_0001",
			Details: fmt.Sprintf("ERROR: %v", err),
			Time:    time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
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
