package middlewares

import (
	"context"
	"net/http"
	"strings"

	apperrors "coolbreez.lk/moderator/internal/errors"
	"coolbreez.lk/moderator/internal/utils"
	"github.com/gin-gonic/gin"
)

type authorizationContextKey string

const AuthorizationContextKey authorizationContextKey = "claims"

func AuthorizationHandler(jwtService *utils.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Authorization")
		if accessToken == "" || !strings.HasPrefix(accessToken, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apperrors.AppStdErrorHandler(
				"missing or invalid authorization header",
				"au_9999",
			))
			return
		}
		token := strings.TrimPrefix(accessToken, "Bearer ")
		extractedData, err := jwtService.VerifyJWTToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apperrors.AppStdErrorHandler(
				err.Error(),
				"au_9999",
			))
			return
		}
		ctx := context.WithValue(c.Request.Context(), AuthorizationContextKey, extractedData)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
