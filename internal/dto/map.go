package dto

import enums "coolbreez.lk/moderator/internal/constants"

type UserSessionIntrospection struct {
	UserID   int64          `json:"id"`
	FullName string         `json:"full_name"`
	Role     enums.UserRole `json:"role"`
}

type UserLoginRequiredFields struct {
	UserID      int64          `json:"id"`
	FullName    string         `json:"full_name"`
	Role        enums.UserRole `json:"role"`
	AccessToken string         `json:"access_token"`
}
