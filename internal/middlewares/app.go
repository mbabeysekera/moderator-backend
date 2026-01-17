package middlewares

import (
	"net/http"

	apperrors "coolbreez.lk/moderator/internal/errors"
	"github.com/gin-gonic/gin"
)

func AppIDExtractHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		appID := c.Request.Header.Get("X-App-Id")
		if appID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apperrors.AppStdErrorHandler(
				"missing or invalid app id",
				"app_9999",
			))
			return
		}
		c.Set("app_id", appID)
		c.Next()
	}
}
