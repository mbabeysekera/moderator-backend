package config

import (
	"fmt"
	"os"
)

func envVariableHandler(key string) (string, error) {
	envVar := os.Getenv(key)
	if envVar == "" {
		return "", fmt.Errorf("environment varibale '%s' is not set", key)
	}
	return envVar, nil
}

func GetServerPort() (string, error) {
	val, err := envVariableHandler("APP_SERVER_PORT")
	if err != nil {
		return "", err
	}
	return val, nil
}

func GetBasePath() (string, error) {
	val, err := envVariableHandler("APP_BASE_PATH")
	if err != nil {
		return "", err
	}
	return val, nil
}
