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

func (r *adminRepo) GetReports(ctx context.Context) ([]review.ReportDetail, error) {
	query := `
	with rp as (select review_id, count(id) as total,
	jsonb_agg(json_build_object('id', id, 'problem', problem, 'reported_at', reported_at)) as problems from reports group by review_id)

	select r.id, r.content, b.name, b.id, u.username, u.id, rp.total, rp.problems from
	rp join reviews as r on rp.review_id = r.id join books as b on r.book_id = b.id join users as u on r.user_id = u.id
	`

	var reports []review.ReportDetail
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rep review.ReportDetail

		err := rows.Scan(&rep.ReviewId, &rep.ReviewContent, &rep.BookName, &rep.BookId, &rep.Username, &rep.UserId, &rep.Total, &rep.Reports)
		if err != nil {
			return nil, err
		}

		reports = append(reports, rep)
	}

	return reports, nil
}

func (r *adminRepo) DeleteFakeReports(ctx context.Context, reviewId int) error {
	res, err := r.db.Exec(ctx, "delete from reports where review_id = $1", reviewId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return review.ErrNotFound
	}

	return nil
}
