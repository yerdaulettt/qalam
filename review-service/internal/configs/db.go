package configs

import (
	"fmt"
	"os"
)

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SslMode  string
}

func NewPostgresConfig() (PostgresConfig, error) {
	cfg := PostgresConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   getEnv("DB_NAME", "review_service"),
		SslMode:  getEnv("DB_SSL", "disable"),
	}

	if cfg.Password == "" || cfg.User == "" || cfg.Host == "" || cfg.Port == "" {
		return PostgresConfig{}, ErrEnv
	}

	return cfg, nil
}

func (p *PostgresConfig) NewDBUrl() string {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", p.User, p.Password, p.Host, p.Port, p.DBName, p.SslMode)

	return dbUrl
}
