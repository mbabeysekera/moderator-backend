package apperrors

import (
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
)

func AppStdErrorHandler(message string, errorID string) *dto.ErrorStdResponse {
	return &dto.ErrorStdResponse{
		Status:  enums.RequestFailed,
		Message: message,
		ErrorID: errorID,
		Time:    time.Now().UTC(),
	}
}
