package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DbUrl         string
	JwtSecret     []byte
	JwtAccessTtl  time.Duration
	JwtRefreshTtl time.Duration
}

func NewConfig() (Config, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "user_service")
	dbSsl := getEnv("DB_SSL", "disable")
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	if dbPassword == "" || dbUser == "" || len(jwtSecret) == 0 {
		return Config{}, ErrEnv
	}

	jwtAccessTtl, err := strconv.Atoi(os.Getenv("JWT_ACCESS_TTL"))
	if err != nil {
		jwtAccessTtl = 15
	}

	jwtRefreshTtl, err := strconv.Atoi(os.Getenv("JWT_REFRESH_TTL"))
	if err != nil {
		jwtRefreshTtl = 24
	}

	cfg := Config{
		DbUrl:         fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSsl),
		JwtSecret:     jwtSecret,
		JwtAccessTtl:  time.Minute * time.Duration(jwtAccessTtl),
		JwtRefreshTtl: time.Hour * time.Duration(jwtRefreshTtl),
	}

	return cfg, nil
}
