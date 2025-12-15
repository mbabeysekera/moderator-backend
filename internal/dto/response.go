package dto

import (
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
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
	AccessToken string              `json:"access_token"`
	Time        time.Time           `json:"time"`
}

type ProductsWithItemsResponse struct {
	All []repositories.ProductsWithItems `json:"all"`
}
