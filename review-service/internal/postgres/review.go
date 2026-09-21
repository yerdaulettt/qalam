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

func (r *reviewRepo) GetReviews(ctx context.Context, bookId int) ([]review.ReviewDetail, error) {
	query := `select r.id, r.content, u.username, r.book_id, coalesce(sum(rl.liked::int), 0) as likes from
	reviews as r join users as u on r.user_id = u.id left join
	review_likes as rl on r.id = rl.review_id where r.book_id = $1 group by r.id, u.id
	`
	var reviews []review.ReviewDetail

	rows, err := r.db.Query(ctx, query, bookId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rev review.ReviewDetail
		if err := rows.Scan(&rev.Id, &rev.Content, &rev.Username, &rev.BookId, &rev.Likes); err != nil {
			return nil, err
		}

		reviews = append(reviews, rev)
	}

	return reviews, nil
}

func (r *reviewRepo) GetMyReviews(ctx context.Context, userId int) ([]review.UserReview, error) {
	query := `
	select r.id, r.content, b.name, b.id, coalesce(sum(rl.liked::int), 0) as likes from
	books as b join reviews as r on b.id = r.book_id left join
	review_likes as rl on r.id = rl.review_id where r.user_id = $1 group by b.id, r.id
	`

	var reviews []review.UserReview
	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rev review.UserReview
		if err := rows.Scan(&rev.Id, &rev.Content, &rev.BookName, &rev.BookId, &rev.Likes); err != nil {
			return nil, err
		}

		reviews = append(reviews, rev)
	}

	return reviews, nil
}

func (r *reviewRepo) GetUserReviews(ctx context.Context, username string) ([]review.UserReview, error) {
	query := `
	select r.id, r.content, b.name, b.id, coalesce(sum(rl.liked::int), 0) as likes from
	(books as b join reviews as r on b.id = r.book_id) join users as u on r.user_id = u.id
	left join review_likes as rl on r.id = rl.review_id where u.username = $1 group by b.id, r.id
	`

	var reviews []review.UserReview
	rows, err := r.db.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rev review.UserReview
		if err := rows.Scan(&rev.Id, &rev.Content, &rev.BookName, &rev.BookId, &rev.Likes); err != nil {
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

func (r *reviewRepo) ReviewLike(ctx context.Context, reviewId, userId int, liked bool) error {
	query := `insert into review_likes (review_id, user_id, liked) values ($1, $2, $3) on conflict
	(review_id, user_id) do update set liked = $3
	`

	_, err := r.db.Exec(ctx, query, reviewId, userId, liked)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == "23503" {
			return review.ErrNotFound
		}

		return err
	}

	return nil
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
	res, err := r.db.Exec(ctx, "delete from books where id = $1", bookId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return review.ErrNotFound
	}

	return nil
}

func (r *reviewRepo) UpdateBook(ctx context.Context, bookId int, name string) error {
	res, err := r.db.Exec(ctx, "update books set name = $1 where id = $2", name, bookId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return review.ErrNotFound
	}

	return nil
}
