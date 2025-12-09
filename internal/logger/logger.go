package logger

import (
	"log/slog"
	"os"

	"coolbreez.lk/moderator/config"
)

func New() *slog.Logger {
	deployEnv, err := config.GetDeployEnv()

	logLevel := slog.LevelInfo
	if err == nil {
		switch deployEnv {
		case "development":
			logLevel = slog.LevelDebug
		case "production":
			logLevel = slog.LevelInfo
		}
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	logger := slog.New(handler).With(
		"app", "coolbreez-moderator",
		"env", deployEnv,
	)

	slog.SetDefault(logger)
	return logger
}
