package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"coolbreez.lk/moderator/internal/dto"
	apperrors "coolbreez.lk/moderator/internal/errors"
	"coolbreez.lk/moderator/internal/services"
	"github.com/gin-gonic/gin"
)

type LoginService interface {
	UserLogin(rc context.Context, loginUser *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
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
	err := c.ShouldBindJSON(&loginUser)
	if err != nil {
		slog.Error("user login parameter validation",
			"err", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)
		c.JSON(http.StatusBadRequest,
			apperrors.AppStdErrorHandler(
				"parameter validation failed",
				"us_0000",
			),
		)
		return
	}
	loginRes, err := lc.service.UserLogin(c.Request.Context(), &loginUser)
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
	c.JSON(http.StatusOK, loginRes)
}
