package configs

import (
	"errors"
	"os"
)

var (
	ErrEnv       = errors.New("Env file incorrect")
	ErrToken     = errors.New("Incorrect token")
	ErrNoToken   = errors.New("No token")
	ErrTokenTime = errors.New("Token expired")
)

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaultValue
}
