package services

import (
	"context"
	"errors"
	"log/slog"

	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/models"
	"coolbreez.lk/moderator/internal/repositories"
	"coolbreez.lk/moderator/internal/utils"
)

type UserServiceImpl struct {
	userRepo *repositories.UserRepository
	jwtUtil  *utils.JWTUtil
}

func NewUserService(repo *repositories.UserRepository, jwtSvc *utils.JWTUtil) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: repo,
		jwtUtil:  jwtSvc,
	}
}

func (us *UserServiceImpl) UserUpdateDetails(rc context.Context,
	userNewDetails *dto.UserUpdateDetails) error {
	user := &models.User{
		ID:       userNewDetails.ID,
		MobileNo: userNewDetails.MobileNo,
		Email:    userNewDetails.Email,
		FullName: userNewDetails.FullName,
	}
	err := us.userRepo.UpdateUserByID(rc, user)
	if err != nil {
		slog.Error("user details update",
			"service", "user",
			"err", err,
			"action", "update",
			"mobile_no", userNewDetails.MobileNo,
		)
		if errors.Is(err, repositories.ErrRowsNotAffected) {
			return ErrUserDetailsUpdate
		}
		return err
	}
	slog.Info("user details updated",
		"service", "user",
		"action", "update",
		"mobile_no", userNewDetails.MobileNo,
	)
	return nil
}
