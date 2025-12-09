package middlewares

import (
	"net/http"

	enums "coolbreez.lk/moderator/internal/constants"
	apperrors "coolbreez.lk/moderator/internal/errors"
	"coolbreez.lk/moderator/internal/utils"
	"github.com/gin-gonic/gin"
)

func CheckRBAC(allowedRoles ...enums.UserRole) gin.HandlerFunc {
	roles := make(map[enums.UserRole]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		roles[role] = struct{}{}
	}
	return func(c *gin.Context) {
		extractedData := c.Request.Context().Value(AuthorizationContextKey)
		if extractedData == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apperrors.AppStdErrorHandler(
				"missing required roles",
				"au_9999",
			))
			return
		}
		val, ok := extractedData.(*utils.JWTExtractedDetails)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apperrors.AppStdErrorHandler(
				"invalid claims type",
				"au_9999",
			))
			return
		}
		_, ok = roles[enums.UserRole(val.UserRole)]
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, apperrors.AppStdErrorHandler(
				"access denied",
				"au_9999",
			))
			return
		}
		c.Next()
	}
}
