package review

import "context"

type reviewRepository interface {
	GetReviews(ctx context.Context, bookId int) ([]Review, error)
	GetMyReviews(ctx context.Context, userId int) ([]UserReview, error)
	GetUserReviews(ctx context.Context, username string) ([]UserReview, error)
	AddReview(ctx context.Context, newReview ReviewReq) error
	ReviewLike(ctx context.Context, reviewId, userId int, liked bool) error
	GetUserId(ctx context.Context, reviewId int) (int, error)
	UpdateReview(ctx context.Context, newReview ReviewUpdate) error
	DeleteReview(ctx context.Context, reviewId, userId int) error
	DeleteBookId(ctx context.Context, bookId int) error
	UpdateBook(ctx context.Context, bookId int, name string) error
	ReportReview(ctx context.Context, reviewId int, problem string) error
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

func (s *ReviewService) GetMyReviews(ctx context.Context, userId int) ([]UserReview, error) {
	r, err := s.repo.GetMyReviews(ctx, userId)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *ReviewService) GetUserReviews(ctx context.Context, username string) ([]UserReview, error) {
	r, err := s.repo.GetUserReviews(ctx, username)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *ReviewService) AddReview(ctx context.Context, newReview ReviewReq) error {
	if newReview.Content == "" || newReview.BookId == 0 || newReview.UserId == 0 {
		return ErrEmpty
	}

	if err := s.repo.AddReview(ctx, newReview); err != nil {
		return err
	}

	return nil
}

func (s *ReviewService) ReviewLike(ctx context.Context, reviewId, userId int, liked bool) error {
	err := s.repo.ReviewLike(ctx, reviewId, userId, liked)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, newReview ReviewUpdate) error {
	correctUserId, err := s.repo.GetUserId(ctx, newReview.ReviewId)
	if err != nil {
		return err
	}

	if newReview.UserId != correctUserId {
		return ErrUnauth
	}

	err = s.repo.UpdateReview(ctx, newReview)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, reviewId, userId int) error {
	correctUserId, err := s.repo.GetUserId(ctx, reviewId)
	if err != nil {
		return err
	}

	if userId != correctUserId {
		return ErrUnauth
	}

	err = s.repo.DeleteReview(ctx, reviewId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReviewService) DeleteBookId(ctx context.Context, bookId int) error {
	if err := s.repo.DeleteBookId(ctx, bookId); err != nil {
		return err
	}

	return nil
}

func (s *ReviewService) UpdateBook(ctx context.Context, bookId int, name string) error {
	if err := s.repo.UpdateBook(ctx, bookId, name); err != nil {
		return err
	}

	return nil
}

func (s *ReviewService) ReportReview(ctx context.Context, reviewId int, problem string) error {
	if err := s.repo.ReportReview(ctx, reviewId, problem); err != nil {
		return err
	}

	return nil
}
