package controllers

import (
	"fmt"
	"net/http"

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
	event.Save()
	c.JSON(http.StatusCreated, gin.H{
		"event":   "created",
		"message": "event was added to the queue",
	})
}
