package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"health":  "up",
		"message": "Moderator is up and running",
	})
}
