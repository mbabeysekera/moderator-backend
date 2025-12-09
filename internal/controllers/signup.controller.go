package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	apperrors "coolbreez.lk/moderator/internal/errors"
	"github.com/gin-gonic/gin"
)

type SignUpService interface {
	UserCreate(rc context.Context, newUser *dto.UserCreateRequest) error
}

type SignUpController struct {
	service SignUpService
}

func NewSignUpController(signupService SignUpService) *SignUpController {
	return &SignUpController{
		service: signupService,
	}
}

func (sc *SignUpController) CreateUser(c *gin.Context) {
	var user dto.UserCreateRequest
	err := c.ShouldBindJSON(&user)
	if err != nil {
		slog.Error("user parameter validation",
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

	err = sc.service.UserCreate(c.Request.Context(), &user)
	if err != nil {
		slog.Error("user create failed",
			"err", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"mobile_no", user.MobileNo,
			"ip", c.ClientIP(),
		)
		c.JSON(http.StatusConflict,
			apperrors.AppStdErrorHandler("create failed", "us_0001"))
		return
	}
	c.JSON(http.StatusCreated, &dto.SuccessStdResponse{
		Status:  enums.RequestSuccess,
		Message: "User Created",
		Details: fmt.Sprintf("User created for Mobile No: %s", user.MobileNo),
		Time:    time.Now(),
	})
}
