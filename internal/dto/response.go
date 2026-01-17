package dto

import (
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/models"
	"coolbreez.lk/moderator/internal/repositories"
)

type ErrorStdResponse struct {
	Status  enums.RequestStatus `json:"status"`
	Message string              `json:"message"`
	ErrorID string              `json:"error_id"`
	Time    time.Time           `json:"time"`
}

type SuccessStdResponse struct {
	Status  enums.RequestStatus `json:"status"`
	Message string              `json:"message"`
	Details string              `json:"details,omitempty"`
	Time    time.Time           `json:"time"`
}

type UserLoginResponse struct {
	Status      enums.RequestStatus `json:"status"`
	UserID      int64               `json:"id"`
	FullName    string              `json:"full_name"`
	Role        enums.UserRole      `json:"role"`
	AccessToken string              `json:"access_token"`
	Time        time.Time           `json:"time"`
	AppID       int64               `json:"app_id"`
}

type SessionIntrospection struct {
	Status   enums.RequestStatus `json:"status"`
	UserID   int64               `json:"id"`
	AppID    int64               `json:"app_id"`
	FullName string              `json:"full_name"`
	Role     enums.UserRole      `json:"role"`
	Time     time.Time           `json:"time"`
}

type ProductsWithItemsResponse struct {
	All   []repositories.ProductWithItems `json:"all"`
	Count int                             `json:"count"`
}

type ProductWithItemsResponse struct {
	Product models.Product `json:"product"`
	Items   []models.Item  `json:"items"`
}
