package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/repositories"
	"coolbreez.lk/moderator/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginServiceImpl struct {
	userRepo   *repositories.UserRepository
	jwtService *utils.JWTUtil
}

func NewLoginService(repo *repositories.UserRepository, jwtSvc *utils.JWTUtil) *LoginServiceImpl {
	return &LoginServiceImpl{
		userRepo:   repo,
		jwtService: jwtSvc,
	}
}

func (ul *LoginServiceImpl) UserLogin(rc context.Context,
	loginUser *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	user, err := ul.userRepo.GetUserByMobileNo(rc, loginUser.MobileNo)
	if err != nil {
		slog.Error("user retrieval error",
			"service", "login",
			"err", err,
			"action", "login",
			"mobile_no", loginUser.MobileNo,
		)
		return nil, err
	}
	if user == nil {
		slog.Info("user does not exists",
			"service", "login",
			"action", "login",
			"mobile_no", loginUser.MobileNo,
		)
		return nil, ErrInvalidUser
	}
	if user.FailedLoginAttempts >= 5 {
		return nil, ErrUserLocked
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginUser.Password))
	if err != nil {
		slog.Error("user password mismatch",
			"service", "login",
			"err", err,
			"action", "validation",
			"mobile_no", loginUser.MobileNo,
		)
		nErr := ul.userRepo.IncrementUserLoginFailuresByID(rc, user.ID)
		if nErr != nil {
			if errors.Is(nErr, repositories.ErrUserNotAffected) {
				return nil, ErrInvalidUser
			}
			return nil, nErr
		}
		return nil, ErrInvalidUser
	}

	accessToke, err := ul.jwtService.GenerateJWTToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	err = ul.userRepo.UpdateSuccessfulLoginByID(rc, user.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotAffected) {
			return nil, ErrInvalidUser
		}
		return nil, err
	}

	slog.Info("user authenticated successfuly",
		"service", "login",
		"action", "login",
		"mobile_no", loginUser.MobileNo,
	)
	return &dto.UserLoginResponse{
		Status:      enums.RequestSuccess,
		AccessToken: accessToke,
		Time:        time.Now().UTC(),
	}, nil
}
