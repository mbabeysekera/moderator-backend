package models

type User struct {
	ID                  int64  `json: "id"`
	MobileNo            string `json: "mobile_no"`
	Email               string `json: "mobile_no"`
	Username            string `json: "username"`
	PasswordHash        string `json: "-"`
	FullName            string `json: "full_name"`
	IsActive            string `json: "is_active"`
	FailedLoginAttempts string `json: "failed_login_attempts"`
	LastLoginAt         string `json: "last_login_at"`
	CreatedAt           string `json: "created_at"`
	UpdatedAt           string `json: "updated_at"`
}
