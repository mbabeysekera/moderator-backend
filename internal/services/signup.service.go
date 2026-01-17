package services

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/models"
	"coolbreez.lk/moderator/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type SignUpServiceImpl struct {
	userRepo *repositories.UserRepository
}

func NewSignUpService(repo *repositories.UserRepository) *SignUpServiceImpl {
	return &SignUpServiceImpl{
		userRepo: repo,
	}
}

func (ss *SignUpServiceImpl) UserCreate(rc context.Context, newUser *dto.UserCreateRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("user password hashing",
			"service", "signup",
			"err", err,
			"action", "generate",
			"mobile_no", newUser.MobileNo,
			"app_id", newUser.AppID,
		)
		return ErrUserDetailsUpdate
	}
	userToCreate := &models.User{
		MobileNo:     newUser.MobileNo,
		Email:        strings.ToLower(strings.TrimSpace(newUser.Email)),
		PasswordHash: string(hashedPassword),
		FullName:     newUser.FullName,
		Role:         enums.RoleUser,
		IsActive:     true,
		AppID:        newUser.AppID,
	}
	err = ss.userRepo.Create(rc, userToCreate)
	if err != nil {
		slog.Error("signup service",
			"service", "signup",
			"err", err,
			"action", "create",
			"mobile_no", newUser.MobileNo,
			"app_id", newUser.AppID,
		)
		if errors.Is(err, repositories.ErrRowsNotAffected) {
			return ErrUserDetailsUpdate
		}
		return err
	}
	slog.Info("user added successfuly",
		"service", "signup",
		"action", "create",
		"mobile_no", newUser.MobileNo,
		"app_id", newUser.AppID,
	)
	return nil
}
