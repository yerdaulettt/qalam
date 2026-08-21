package postgres

import (
	"context"

	"review-service/internal/review"

	"github.com/jackc/pgx/v5/pgxpool"
)

type reviewRepo struct {
	db *pgxpool.Pool
}

func NewReviewRepo(db *pgxpool.Pool) *reviewRepo {
	return &reviewRepo{db: db}
}

func (r *reviewRepo) DeleteBookId(ctx context.Context, bookId int) error {
	res, err := r.db.Exec(ctx, "delete from books_id where id = $1", bookId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return review.ErrNotFound
	}

	return nil
}
