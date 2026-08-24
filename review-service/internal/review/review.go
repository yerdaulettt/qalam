package review

import "context"

type reviewRepository interface {
	GetReviews(ctx context.Context, bookId int) ([]Review, error)
	AddReview(ctx context.Context, newReview ReviewReq) (Review, error)
	GetUserId(ctx context.Context, reviewId int) (int, error)
	UpdateReview(ctx context.Context, newReview ReviewUpdate) (Review, error)
	DeleteReview(ctx context.Context, reviewId, userId int) (Review, error)
	DeleteBookId(ctx context.Context, bookId int) error
}

type ReviewService struct {
	repo reviewRepository
}

func NewReviewService(r reviewRepository) *ReviewService {
	return &ReviewService{repo: r}
}

func (s *ReviewService) GetReviews(ctx context.Context, bookId int) ([]Review, error) {
	r, err := s.repo.GetReviews(ctx, bookId)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *ReviewService) AddReview(ctx context.Context, newReview ReviewReq) (Review, error) {
	if newReview.Content == "" {
		return Review{}, ErrEmpty
	}

	r, err := s.repo.AddReview(ctx, newReview)
	if err != nil {
		return Review{}, err
	}

	return r, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, newReview ReviewUpdate) (Review, error) {
	correctUserId, err := s.repo.GetUserId(ctx, newReview.ReviewId)
	if err != nil {
		return Review{}, err
	}

	if newReview.UserId != correctUserId {
		return Review{}, ErrUnauth
	}

	rev, err := s.repo.UpdateReview(ctx, newReview)
	if err != nil {
		return Review{}, err
	}

	return rev, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, reviewId, userId int) (Review, error) {
	correctUserId, err := s.repo.GetUserId(ctx, reviewId)
	if err != nil {
		return Review{}, err
	}

	if userId != correctUserId {
		return Review{}, ErrUnauth
	}

	rev, err := s.repo.DeleteReview(ctx, reviewId, userId)
	if err != nil {
		return Review{}, err
	}

	return rev, nil
}

func (s *ReviewService) DeleteBookId(ctx context.Context, bookId int) error {
	if err := s.repo.DeleteBookId(ctx, bookId); err != nil {
		return err
	}

	return nil
}
