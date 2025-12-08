package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(rc context.Context, newUser *dto.UserCreateRequest) error
}

type UserController struct {
	service UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{
		service: userService,
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	log.Printf("user create request")
	var user dto.UserCreateRequest
	c.ShouldBindJSON(&user)
	err := uc.service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorStdResponse{
			Status:  enums.RequestFailed,
			Message: fmt.Sprintf("error from user.controller.create[DATA]: %v", c.Request.Body),
			ErrorID: "us_0000",
			Details: fmt.Sprintf("ERROR: %v", err),
			Time:    time.Now(),
		})
		return
	}
	c.JSON(http.StatusCreated, &dto.CreateStdResponse{
		Status:  enums.RequestSuccess,
		Message: "User Created",
		Details: "",
		Time:    time.Now(),
	})
}
