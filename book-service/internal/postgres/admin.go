package postgres

import (
	"context"
	"errors"

	"book-service/internal/book"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type adminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) NewGenre(ctx context.Context, name string) (book.Genre, error) {
	query := "insert into genres (name) values ($1) returning id, name"

	var g book.Genre
	if err := r.db.QueryRow(ctx, query, name).Scan(&g.Id, &g.Name); err != nil {
		return book.Genre{}, err
	}

	return g, nil
}

func (r *adminRepo) UpdateGenre(ctx context.Context, genreId int, name string) (book.Genre, error) {
	query := "update genres set name = $1 where id = $2 returning id, name"

	var g book.Genre
	if err := r.db.QueryRow(ctx, query, name, genreId).Scan(&g.Id, &g.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return book.Genre{}, book.ErrNotFound
		}

		return book.Genre{}, err
	}

	return g, nil
}

func (r *adminRepo) DeleteGenre(ctx context.Context, genreId int) (book.Genre, error) {
	query := "delete from genres where id = $1 returning id, name"

	var g book.Genre
	if err := r.db.QueryRow(ctx, query, genreId).Scan(&g.Id, &g.Name); err != nil {
		return book.Genre{}, err
	}

	return g, nil
}
