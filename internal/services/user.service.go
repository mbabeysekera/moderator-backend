package services

import (
	"context"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/models"
	"coolbreez.lk/moderator/internal/repositories"
)

type UserServiceImpl struct {
	userRepo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: repo,
	}
}

func (us *UserServiceImpl) CreateUser(rc context.Context, newUser *dto.UserCreateRequest) error {
	userToCreate := &models.User{
		MobileNo:     newUser.MobileNo,
		Email:        newUser.Email,
		PasswordHash: newUser.Password,
		FullName:     newUser.FullName,
		Role:         enums.RoleUser,
		IsActive:     true,
	}
	err := us.userRepo.Create(rc, userToCreate)
	if err != nil {
		return err
	}
	return nil
}
