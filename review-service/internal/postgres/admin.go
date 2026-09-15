package postgres

import (
	"context"
	"errors"

	"review-service/internal/review"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type adminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) GetAllReviews(ctx context.Context) ([]review.ReviewDetailAdmin, error) {
	query := `
	select r.id, r.content, b.name, b.id, u.username, u.id from
	reviews as r join books as b on r.book_id = b.id
	join users as u on r.user_id = u.id
	`

	var reviews []review.ReviewDetailAdmin

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rev review.ReviewDetailAdmin

		err := rows.Scan(&rev.Id, &rev.Content, &rev.BookName, &rev.BookId, &rev.Username, &rev.UserId)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, rev)
	}

	return reviews, nil
}

func (r *adminRepo) DeleteReview(ctx context.Context, reviewId int) (review.Review, error) {
	query := "delete from reviews where id = $1 returning id, content, user_id, book_id"
	var rev review.Review

	err := r.db.QueryRow(ctx, query, reviewId).Scan(&rev.Id, &rev.Content, &rev.UserId, &rev.BookId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return review.Review{}, review.ErrNotFound
		}

		return review.Review{}, err
	}

	return rev, nil
}
