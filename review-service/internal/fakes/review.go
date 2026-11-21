package fakes

import (
	"context"

	"review-service/internal/review"
)

type fakeReview struct {
	Reviews map[int]review.Review
	Reports []string
	Id      int
}

func NewFakeReview() *fakeReview {
	return &fakeReview{Reviews: make(map[int]review.Review)}
}

func (f *fakeReview) GetReviews(ctx context.Context, bookId int) ([]review.Review, error) {
	return nil, nil
}

func (f *fakeReview) GetMyReviews(ctx context.Context, userId int) ([]review.UserReview, error) {
	return nil, nil
}

func (f *fakeReview) GetUserReviews(ctx context.Context, username string) ([]review.UserReview, error) {
	return nil, nil
}

func (f *fakeReview) AddReview(ctx context.Context, newReview review.ReviewReq) error {
	f.Id++
	f.Reviews[f.Id] = review.Review{Id: f.Id, Content: newReview.Content, UserId: newReview.UserId, BookId: newReview.BookId}

	return nil
}

func (f *fakeReview) ReviewLike(ctx context.Context, reviewId, userId int, liked bool) error {
	r, ok := f.Reviews[reviewId]
	if !ok {
		return review.ErrNotFound
	}

	if liked {
		r.Likes++
		f.Reviews[reviewId] = r
	}

	return nil
}

func (f *fakeReview) GetUserId(ctx context.Context, reviewId int) (int, error) {
	u, ok := f.Reviews[reviewId]
	if !ok {
		return 0, review.ErrNotFound
	}

	return u.UserId, nil
}

func (f *fakeReview) UpdateReview(ctx context.Context, newReview review.ReviewUpdate) error {
	r, ok := f.Reviews[newReview.ReviewId]
	if !ok {
		return review.ErrNotFound
	}

	r.Content = newReview.Content
	f.Reviews[r.Id] = r

	return nil
}

func (f *fakeReview) DeleteReview(ctx context.Context, reviewId, userId int) error {
	_, ok := f.Reviews[reviewId]
	if !ok {
		return review.ErrNotFound
	}

	delete(f.Reviews, reviewId)

	return nil
}

func (f *fakeReview) DeleteBookId(ctx context.Context, bookId int) error {
	return nil
}

func (f *fakeReview) UpdateBook(ctx context.Context, bookId int, name string) error {
	return nil
}

func (f *fakeReview) ReportReview(ctx context.Context, reviewId int, problem string) error {
	_, ok := f.Reviews[reviewId]
	if !ok {
		return review.ErrNotFound
	}

	f.Reports = append(f.Reports, problem)

	return nil
}
