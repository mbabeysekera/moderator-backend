package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/middlewares"
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
	userID := rc.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserID
	userRole := rc.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserRole
	appID := rc.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).AppID

	if enums.UserRole(userRole) != enums.RoleAdmin &&
		appID != userNewDetails.AppID &&
		userID != userNewDetails.ID {
		slog.Error("user details update",
			"service", "user",
			"err", ErrInvalidUser.Error(),
			"action", "update",
			"user_id", userID,
		)
		return ErrInvalidUser
	}

	user, err := us.userRepo.GetUserByID(rc, userNewDetails.ID)
	if err != nil {
		slog.Error("user details update",
			"service", "user",
			"err", err,
			"action", "update",
			"user_id", userNewDetails.ID,
		)
		return err
	}
	if user == nil {
		slog.Error("user details update",
			"service", "user",
			"err", err,
			"action", "update",
			"user_id", userNewDetails.ID,
		)
		return ErrInvalidUser
	}

	if userNewDetails.MobileNo != "" && userNewDetails.MobileNo != user.MobileNo {
		existingUser, err := us.userRepo.GetUserByMobileNo(rc, userNewDetails.MobileNo, user.AppID)
		if err != nil {
			return err
		}
		if existingUser != nil && existingUser.ID != user.ID {
			slog.Error("user details update",
				"service", "user",
				"err", "mobile number already exists",
				"action", "update",
				"user_id", userNewDetails.ID,
			)
			return errors.New("mobile number already exists")
		}
	}

	if userNewDetails.Email != "" && userNewDetails.Email != user.Email {
		existingUser, err := us.userRepo.GetUserByEmail(rc, userNewDetails.Email, user.AppID)
		if err != nil {
			return err
		}
		if existingUser != nil && existingUser.ID != user.ID {
			slog.Error("user details update",
				"service", "user",
				"err", "email already exists",
				"action", "update",
				"user_id", userNewDetails.ID,
			)
			return errors.New("email already exists")
		}
	}

	user.MobileNo = userNewDetails.MobileNo
	user.Email = userNewDetails.Email
	user.FullName = userNewDetails.FullName
	user.UpdatedAt = time.Now().UTC()

	err = us.userRepo.UpdateUserByID(rc, user)
	if err != nil {
		slog.Error("user details update",
			"service", "user",
			"err", err,
			"action", "update",
			"user_id", userID,
		)
		if errors.Is(err, repositories.ErrRowsNotAffected) {
			return ErrUserDetailsUpdate
		}
		return err
	}
	slog.Info("user details updated",
		"service", "user",
		"action", "update",
		"user_id", userID,
	)
	return nil
}

func (us *UserServiceImpl) GetUserByAccessToken(rc context.Context) (*dto.UserSessionIntrospection,
	error) {
	userID := rc.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserID
	user, err := us.userRepo.GetUserByID(rc, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidUser
	}
	slog.Info("user details fetch",
		"service", "user",
		"action", "fetch",
		"user_id", userID,
	)
	return &dto.UserSessionIntrospection{
		UserID:   user.ID,
		FullName: user.FullName,
		Role:     user.Role,
		AppID:    user.AppID,
	}, nil
}
