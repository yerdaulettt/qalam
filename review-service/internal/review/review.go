package review

import "context"

type reviewRepository interface {
	DeleteBookId(ctx context.Context, bookId int) error
}

type ReviewService struct {
	repo reviewRepository
}

func NewReviewService(r reviewRepository) *ReviewService {
	return &ReviewService{repo: r}
}

func (s *ReviewService) DeleteBookId(ctx context.Context, bookId int) error {
	if err := s.repo.DeleteBookId(ctx, bookId); err != nil {
		return err
	}

	return nil
}
