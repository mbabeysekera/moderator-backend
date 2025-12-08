package models

import (
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
)

type Event struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Target      string          `json:"target"`
	Command     string          `json:"command"`
	Status      enums.EventCode `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	ProcessedAt *time.Time      `json:"processed_at,omitempty"`
}
