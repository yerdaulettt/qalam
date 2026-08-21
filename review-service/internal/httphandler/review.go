package httphandler

import (
	"review-service/internal/review"
)

type reviewHandler struct {
	service *review.ReviewService
}

func NewReviewHandler(s *review.ReviewService) *reviewHandler {
	return &reviewHandler{service: s}
}
