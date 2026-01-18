package middlewares

import (
	"net/http"
	"strconv"

	apperrors "coolbreez.lk/moderator/internal/errors"
	"github.com/gin-gonic/gin"
)

func AppIDExtractHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		appIDStr := c.Request.Header.Get("X-App-Id")
		if appIDStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apperrors.AppStdErrorHandler(
				"missing or invalid app id",
				"app_9999",
			))
			return
		}

		appID, err := strconv.ParseInt(appIDStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, apperrors.AppStdErrorHandler(
				"invalid app id format",
				"app_9998",
			))
			return
		}

		c.Set("app_id", appID)
		c.Next()
	}
}
