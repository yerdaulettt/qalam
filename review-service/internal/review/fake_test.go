package review

import "context"

type fakeReview struct {
	reviews map[int]Review
	id      int
	err     []error
}

func newFakeReview() *fakeReview {
	return &fakeReview{reviews: make(map[int]Review), err: make([]error, 2)}
}

func (f *fakeReview) GetReviews(ctx context.Context, bookId int) ([]Review, error) {
	return nil, nil
}

func (f *fakeReview) GetMyReviews(ctx context.Context, userId int) ([]UserReview, error) {
	return nil, nil
}

func (f *fakeReview) GetUserReviews(ctx context.Context, username string) ([]UserReview, error) {
	return nil, nil
}

func (f *fakeReview) AddReview(ctx context.Context, newReview ReviewReq) error {
	if f.err[0] != nil {
		return f.err[0]
	}

	f.id++
	f.reviews[f.id] = Review{Id: f.id, Content: newReview.Content, UserId: newReview.UserId, BookId: newReview.BookId}

	return nil
}

func (f *fakeReview) ReviewLike(ctx context.Context, reviewId, userId int, liked bool) error {
	return nil
}

func (f *fakeReview) GetUserId(ctx context.Context, reviewId int) (int, error) {
	if f.err[1] != nil {
		return 0, f.err[1]
	}

	u, ok := f.reviews[reviewId]
	if !ok {
		return 0, ErrNotFound
	}

	return u.UserId, nil
}

func (f *fakeReview) UpdateReview(ctx context.Context, newReview ReviewUpdate) error {
	return nil
}

func (f *fakeReview) DeleteReview(ctx context.Context, reviewId, userId int) error {
	if f.err[0] != nil {
		return f.err[0]
	}

	_, ok := f.reviews[reviewId]
	if !ok {
		return ErrNotFound
	}

	delete(f.reviews, reviewId)

	return nil
}

func (f *fakeReview) DeleteBookId(ctx context.Context, bookId int) error {
	return nil
}

func (f *fakeReview) UpdateBook(ctx context.Context, bookId int, name string) error {
	return nil
}

func (f *fakeReview) ReportReview(ctx context.Context, reviewId int, problem string) error {
	return nil
}
