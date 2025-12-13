package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/models"
	"coolbreez.lk/moderator/internal/repositories"
	"coolbreez.lk/moderator/internal/utils"
)

var ErrInvalidUser = errors.New("invalid credentials")
var ErrUserLocked = errors.New("user locked")
var ErrUserDetailsUpdate = errors.New("user not affected")

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

// func (us *UserServiceImpl) UserCreate(rc context.Context, newUser *dto.UserCreateRequest) error {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		slog.Error("user password hashing",
// 			"service", "user",
// 			"err", err,
// 			"action", "generate",
// 			"mobile_no", newUser.MobileNo,
// 		)
// 		return ErrUserDetailsUpdate
// 	}
// 	userToCreate := &models.User{
// 		MobileNo:     newUser.MobileNo,
// 		Email:        strings.ToLower(strings.TrimSpace(newUser.Email)),
// 		PasswordHash: string(hashedPassword),
// 		FullName:     newUser.FullName,
// 		Role:         enums.RoleUser,
// 		IsActive:     true,
// 	}
// 	err = us.userRepo.Create(rc, userToCreate)
// 	if err != nil {
// 		slog.Error("user service",
// 			"service", "user",
// 			"err", err,
// 			"action", "create",
// 			"mobile_no", newUser.MobileNo,
// 		)
// 		if errors.Is(err, repositories.ErrUserNotAffected) {
// 			return ErrUserDetailsUpdate
// 		}
// 		return err
// 	}
// 	slog.Info("user added successfuly",
// 		"service", "user",
// 		"action", "create",
// 		"mobile_no", newUser.MobileNo,
// 	)
// 	return nil
// }

// func (us *UserServiceImpl) UserLogin(rc context.Context,
// 	loginUser *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
// 	user, err := us.userRepo.GetUserByMobileNo(rc, loginUser.MobileNo)
// 	if err != nil {
// 		slog.Error("user retrieval error",
// 			"service", "user",
// 			"err", err,
// 			"action", "login",
// 			"mobile_no", loginUser.MobileNo,
// 		)
// 		return nil, err
// 	}
// 	if user == nil {
// 		slog.Info("user does not exists",
// 			"service", "user",
// 			"action", "login",
// 			"mobile_no", loginUser.MobileNo,
// 		)
// 		return nil, ErrInvalidUser
// 	}
// 	if user.FailedLoginAttempts >= 5 {
// 		return nil, ErrUserLocked
// 	}
// 	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginUser.Password))
// 	if err != nil {
// 		slog.Error("user password mismatch",
// 			"service", "user",
// 			"err", err,
// 			"action", "validation",
// 			"mobile_no", loginUser.MobileNo,
// 		)
// 		nErr := us.userRepo.IncrementUserLoginFailuresByID(rc, user.ID)
// 		if nErr != nil {
// 			if errors.Is(nErr, repositories.ErrUserNotAffected) {
// 				return nil, ErrInvalidUser
// 			}
// 			return nil, nErr
// 		}
// 		return nil, ErrInvalidUser
// 	}

// 	accessToke, err := us.jwtUtil.GenerateJWTToken(user.ID, user.Role)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = us.userRepo.UpdateSuccessfulLoginByID(rc, user.ID)
// 	if err != nil {
// 		if errors.Is(err, repositories.ErrUserNotAffected) {
// 			return nil, ErrInvalidUser
// 		}
// 		return nil, err
// 	}

// 	slog.Info("user authenticated successfuly",
// 		"service", "user",
// 		"action", "login",
// 		"mobile_no", loginUser.MobileNo,
// 	)
// 	return &dto.UserLoginResponse{
// 		Status:      enums.RequestSuccess,
// 		AccessToken: accessToke,
// 		Time:        time.Now().UTC(),
// 	}, nil
// }

func (us *UserServiceImpl) UserUpdateDetails(rc context.Context,
	userNewDetails *dto.UserUpdateDetails) (*dto.SuccessStdResponse, error) {
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
		if errors.Is(err, repositories.ErrUserNotAffected) {
			return nil, ErrUserDetailsUpdate
		}
		return nil, err
	}
	slog.Info("user details updated",
		"service", "user",
		"action", "update",
		"mobile_no", userNewDetails.MobileNo,
	)
	return &dto.SuccessStdResponse{
		Status:  enums.RequestSuccess,
		Message: "refresh page to update details",
		Details: "user details updated",
		Time:    time.Now().UTC(),
	}, nil
}
