package dto

import enums "coolbreez.lk/moderator/internal/constants"

type UserSessionIntrospection struct {
	UserID int64          `json:"id"`
	Role   enums.UserRole `json:"role"`
}
