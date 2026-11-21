package postgres

import (
	"context"
	"errors"

	"user-service/internal/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetUsernameStatus(ctx context.Context, username string) (bool, error) {
	var id int

	err := r.db.QueryRow(ctx, "select id from users where username = $1", username).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		}

		return false, err
	}

	if id != 0 {
		return false, nil
	}

	return true, nil
}

func (r *userRepo) GetUser(ctx context.Context, username string) (user.UserVerify, error) {
	var u user.UserVerify

	err := r.db.QueryRow(ctx, "select id, role, hash from users where username = $1", username).Scan(&u.Id, &u.Role, &u.Hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.UserVerify{}, user.ErrNotFound
		}

		return user.UserVerify{}, err
	}

	return u, nil
}

func (r *userRepo) GetMyProfile(ctx context.Context, userId int) (user.User, error) {
	query := "select id, name, username, role from users where id = $1"
	var u user.User

	err := r.db.QueryRow(ctx, query, userId).Scan(&u.Id, &u.Name, &u.Username, &u.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, user.ErrNotFound
		}

		return user.User{}, err
	}

	return u, nil
}

func (r *userRepo) Register(ctx context.Context, newUser user.RegisterReq) (user.User, error) {
	query := "insert into users (name, username, role, hash) values ($1, $2, $3, $4) returning id, name, username, role"
	var u user.User

	err := r.db.QueryRow(ctx, query, newUser.Name, newUser.Username, newUser.Role, newUser.Password).Scan(&u.Id, &u.Name, &u.Username, &u.Role)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}
