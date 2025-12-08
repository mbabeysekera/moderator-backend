package dto

import (
	enums "coolbreez.lk/moderator/internal/constants"
)

type EventRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Target      string          `json:"target"`
	Command     string          `json:"command"`
	Status      enums.EventCode `json:"status"`
}

type UserCreateRequest struct {
	MobileNo string `json:"mobile_no"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}
