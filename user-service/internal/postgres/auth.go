package postgres

import (
	"context"
	"errors"

	"user-service/internal/auth"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *authRepo {
	return &authRepo{db: db}
}

func (r *authRepo) GetUsername(ctx context.Context, username string) (string, error) {
	var uname string

	err := r.db.QueryRow(ctx, "select id from users where username = $1", username).Scan(&uname)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return uname, nil
}

func (r *authRepo) GetUser(ctx context.Context, username string) (auth.UserVerify, error) {
	var user auth.UserVerify

	err := r.db.QueryRow(ctx, "select id, role, hash from users where username = $1", username).Scan(&user.Id, &user.Role, &user.Hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.UserVerify{}, auth.ErrNotFound
		}

		return auth.UserVerify{}, err
	}

	return user, nil
}

func (r *authRepo) Register(ctx context.Context, u auth.RegisterReq) (auth.User, error) {
	query := "insert into users (name, username, role, hash) values ($1, $2, $3, $4) returning id, name, username, role"
	var user auth.User

	err := r.db.QueryRow(ctx, query, u.Name, u.Username, u.Role, u.Password).Scan(&user.Id, &user.Name, &user.Username, &user.Role)
	if err != nil {
		return auth.User{}, err
	}

	return user, nil
}
