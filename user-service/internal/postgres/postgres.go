package postgres

import (
	"context"
	"fmt"

	"user-service/internal/configs"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, cfg *configs.PostgresConfig) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SslMode)

	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
