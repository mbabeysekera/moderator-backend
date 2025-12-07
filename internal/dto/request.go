package dto

import (
	"coolbreez.lk/moderator/internal/constants"
)

type EventRequest struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Target      string              `json:"target"`
	Command     string              `json:"command"`
	Status      constants.EventCode `json:"status"`
}
