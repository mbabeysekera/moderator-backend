package models

import (
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
)

type User struct {
	ID                  int64          `json:"id"`
	MobileNo            string         `json:"mobile_no"`
	Email               string         `json:"email,omitempty"`
	PasswordHash        string         `json:"-"`
	FullName            string         `json:"full_name"`
	Role                enums.UserRole `json:"role"`
	IsActive            bool           `json:"is_active"`
	FailedLoginAttempts int            `json:"failed_login_attempts"`
	LastLoginAt         *time.Time     `json:"last_login_at"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}
