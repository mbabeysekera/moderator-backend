package controllers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	apperrors "coolbreez.lk/moderator/internal/errors"
	"coolbreez.lk/moderator/internal/services"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	UserUpdateDetails(rc context.Context, userNewDetails *dto.UserUpdateDetails) error
	GetUserByID(rc context.Context) (*dto.UserSessionIntrospection, error)
}

type UserController struct {
	service UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{
		service: userService,
	}
}

func (uc *UserController) UserDetailsUpdate(c *gin.Context) {
	var userNewDetails dto.UserUpdateDetails
	err := c.ShouldBindJSON(&userNewDetails)
	if err != nil {
		slog.Error("user update details",
			"err", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)
		c.JSON(http.StatusBadRequest, apperrors.AppStdErrorHandler(
			"parameter validation failed",
			"us_0000",
		))
		return
	}
	err = uc.service.UserUpdateDetails(c.Request.Context(), &userNewDetails)
	if err != nil {
		slog.Error("user details update",
			"err", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"mobile_no", userNewDetails.MobileNo,
			"ip", c.ClientIP(),
		)
		if errors.Is(err, services.ErrUserDetailsUpdate) {
			c.JSON(http.StatusBadRequest, apperrors.AppStdErrorHandler(
				"User detais not updated",
				"us_0001",
			))
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"us_0002",
			),
		)
		return
	}
	c.JSON(http.StatusOK, &dto.SuccessStdResponse{
		Status:  enums.RequestSuccess,
		Message: "refresh page to update details",
		Details: "user details updated",
		Time:    time.Now().UTC(),
	})
}

func (uc *UserController) UserSessionIntrospection(c *gin.Context) {
	userSessionIntros, err := uc.service.GetUserByID(c.Request.Context())
	if err != nil {
		slog.Error("user details update",
			"err", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"user_id", "",
			"ip", c.ClientIP(),
		)
		if errors.Is(err, services.ErrInvalidUser) {
			c.JSON(http.StatusBadRequest,
				apperrors.AppStdErrorHandler(
					services.ErrInvalidUser.Error(),
					"us_0000",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"us_0001",
			),
		)
		return
	}
	c.JSON(http.StatusOK, &dto.SessionIntrospection{
		Status:   enums.RequestSuccess,
		UserID:   userSessionIntros.UserID,
		FullName: userSessionIntros.FullName,
		Role:     userSessionIntros.Role,
		Time:     time.Now().UTC(),
	})
}
