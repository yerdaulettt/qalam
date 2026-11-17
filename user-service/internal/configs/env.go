package configs

import (
	"errors"
	"os"
)

var (
	ErrEnv = errors.New("Env file incorrect")
)

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaultValue
}
