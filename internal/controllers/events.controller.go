package controllers

import (
	"fmt"
	"net/http"

	"coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context) {
	var event models.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		// log.Fatalf("error from event.controller.create[DATA]: %v [ERROR]: %v", c.Request.Body, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Something is not right.",
			"message": fmt.Sprintf("error from event.controller.create[DATA]: %v", c.Request.Body),
		})
		return
	}
	event.Save(c.Request.Context())
	c.JSON(http.StatusCreated, gin.H{
		"event":   "created",
		"message": "event was added to the queue",
	})
}

func GetEventsByStatus(c *gin.Context) {
	status := c.Param("status")
	var event models.Event
	events, err := event.GetEventsByStatus(c.Request.Context(), constants.EventCode(status))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Something is not right.",
			"message": fmt.Sprintf("error from event.controller.getEventsByStatus: %v", c.Request.Body),
		})
	}
	c.JSON(http.StatusOK, events)
}
