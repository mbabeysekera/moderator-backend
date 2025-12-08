package dto

import (
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
)

type ErrorStdResponse struct {
	Status  enums.RequestStatus `json:"status"`
	Message string              `json:"message"`
	ErrorID string              `json:"error_id"`
	Details string              `json:"details,omitempty"`
	Time    time.Time           `json:"time"`
}

type CreateStdResponse struct {
	Status  enums.RequestStatus `json:"status"`
	Message string              `json:"message"`
	Details string              `json:"details,omitempty"`
	Time    time.Time           `json:"time"`
}
