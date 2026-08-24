package postgres

import (
	"context"
	"errors"

	"review-service/internal/review"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type reviewRepo struct {
	db *pgxpool.Pool
}

func NewReviewRepo(db *pgxpool.Pool) *reviewRepo {
	return &reviewRepo{db: db}
}

func (r *reviewRepo) GetReviews(ctx context.Context, bookId int) ([]review.Review, error) {
	query := `select id, content, user_id, book_id from reviews where book_id = $1`
	var reviews []review.Review

	rows, err := r.db.Query(ctx, query, bookId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rev review.Review
		if err := rows.Scan(&rev.Id, &rev.Content, &rev.UserId, &rev.BookId); err != nil {
			return nil, err
		}

		reviews = append(reviews, rev)
	}

	return reviews, nil
}

func (r *reviewRepo) AddReview(ctx context.Context, newReview review.ReviewReq) (review.Review, error) {
	query := `insert into reviews (content, user_id, book_id) values ($1, $2, $3) returning id, content, user_id, book_id`

	var rev review.Review
	err := r.db.QueryRow(ctx, query, newReview.Content, newReview.UserId, newReview.BookId).Scan(&rev.Id, &rev.Content, &rev.UserId, &rev.BookId)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "23503" {
				return review.Review{}, review.ErrNotFound
			}
		}

		return review.Review{}, err
	}

	return rev, nil
}

func (r *reviewRepo) GetUserId(ctx context.Context, reviewId int) (int, error) {
	var userId int

	err := r.db.QueryRow(ctx, "select user_id from reviews where id = $1", reviewId).Scan(&userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, review.ErrNotFound
		}

		return 0, err
	}

	return userId, nil
}

func (r *reviewRepo) UpdateReview(ctx context.Context, newReview review.ReviewUpdate) (review.Review, error) {
	query := `update reviews set content = $1 where id = $2 and user_id = $3 returning id, content, user_id, book_id`

	var rev review.Review
	err := r.db.QueryRow(ctx, query, newReview.Content, newReview.ReviewId, newReview.UserId).Scan(&rev.Id, &rev.Content, &rev.UserId, &rev.BookId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return review.Review{}, review.ErrNotFound
		}

		return review.Review{}, err
	}

	return rev, nil
}

func (r *reviewRepo) DeleteReview(ctx context.Context, reviewId, userId int) (review.Review, error) {
	query := `delete from reviews where id = $1 and user_id = $2 returning id, content, user_id, book_id`

	var rev review.Review
	err := r.db.QueryRow(ctx, query, reviewId, userId).Scan(&rev.Id, &rev.Content, &rev.UserId, &rev.BookId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return review.Review{}, review.ErrNotFound
		}

		return review.Review{}, err
	}

	return rev, nil
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
