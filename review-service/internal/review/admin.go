package review

import "context"

type adminRepository interface {
	GetAllReviews(ctx context.Context) ([]ReviewDetailAdmin, error)
	DeleteReview(ctx context.Context, reviewId int) (Review, error)
}

type AdminService struct {
	repo adminRepository
}

func NewAdminService(r adminRepository) *AdminService {
	return &AdminService{repo: r}
}

func (s *AdminService) GetAllReviews(ctx context.Context) ([]ReviewDetailAdmin, error) {
	reviews, err := s.repo.GetAllReviews(ctx)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *AdminService) DeleteReview(ctx context.Context, reviewId int) (Review, error) {
	rev, err := s.repo.DeleteReview(ctx, reviewId)
	if err != nil {
		return Review{}, err
	}

	return rev, nil
}
