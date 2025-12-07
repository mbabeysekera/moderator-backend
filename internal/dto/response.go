package dto

import (
	"time"

	"coolbreez.lk/moderator/internal/constants"
)

type ErrorStdResponse struct {
	Status  constants.RequestStatus `json: "status"`
	Message string                  `json: "message"`
	ErrorID string                  `json: "error_id"`
	Details string                  `json: "details,omitempty"`
	Time    time.Time               `json: "time"`
}

type EventResponse struct {
	ID          int64               `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Target      string              `json:"target"`
	Command     string              `json:"command"`
	Status      constants.EventCode `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	ProcessedAt *time.Time          `json:"processed_at,omitempty"`
}
