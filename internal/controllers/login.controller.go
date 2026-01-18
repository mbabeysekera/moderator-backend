package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	apperrors "coolbreez.lk/moderator/internal/errors"
	"coolbreez.lk/moderator/internal/services"
	"github.com/gin-gonic/gin"
)

type LoginService interface {
	UserLogin(rc context.Context,
		loginUser *dto.UserLoginRequest, appID int64) (*dto.UserLoginRequiredFields, error)
}

type LoginController struct {
	service LoginService
}

func NewLoginController(loginService LoginService) *LoginController {
	return &LoginController{
		service: loginService,
	}
}

func (lc *LoginController) Login(c *gin.Context) {
	var loginUser dto.UserLoginRequest
	appID := c.GetInt64("app_id")
	err := c.ShouldBindJSON(&loginUser)
	if err != nil {
		slog.Error("user login parameter validation",
			"err", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"app_id", appID,
		)
		c.JSON(http.StatusBadRequest,
			apperrors.AppStdErrorHandler(
				"parameter validation failed",
				"us_0000",
			),
		)
		return
	}
	loginRequiredFields, err := lc.service.UserLogin(c.Request.Context(), &loginUser, appID)
	if err != nil {
		slog.Error("user login unsuccessful",
			"err", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"mobile_no", loginUser.MobileNo,
			"ip", c.ClientIP(),
		)
		if errors.Is(err, services.ErrInvalidUser) {
			c.JSON(http.StatusUnauthorized,
				apperrors.AppStdErrorHandler(
					fmt.Sprintf("user login failed for MobileNo: %s", loginUser.MobileNo),
					"us_0001",
				),
			)
			return
		}
		if errors.Is(err, services.ErrUserLocked) {
			c.JSON(http.StatusUnauthorized,
				apperrors.AppStdErrorHandler(
					fmt.Sprintf("user locked for MobileNo: %s, please contact admin via: %s",
						loginUser.MobileNo,
						"moderator.lk@gmail.com",
					),
					"us_0002",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"us_0003",
			),
		)
		return
	}
	c.JSON(http.StatusOK, &dto.UserLoginResponse{
		Status:      enums.RequestSuccess,
		AccessToken: loginRequiredFields.AccessToken,
		UserID:      loginRequiredFields.UserID,
		FullName:    loginRequiredFields.FullName,
		Role:        loginRequiredFields.Role,
		AppID:       loginRequiredFields.AppID,
		Time:        time.Now().UTC(),
	})
}
