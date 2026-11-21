package review_test

import (
	"errors"
	"slices"
	"testing"

	"review-service/internal/fakes"
	"review-service/internal/review"
)

func TestAddReview(t *testing.T) {
	repo := fakes.NewFakeReview()
	service := review.NewReviewService(repo)

	tests := []struct {
		newReview review.ReviewReq
		err       error
	}{
		{newReview: review.ReviewReq{Content: ""}, err: review.ErrEmpty},
		{newReview: review.ReviewReq{Content: "Test"}, err: review.ErrEmpty},
		{newReview: review.ReviewReq{Content: "Test", BookId: 1}, err: review.ErrEmpty},
		{newReview: review.ReviewReq{Content: "Test", UserId: 1}, err: review.ErrEmpty},
		{newReview: review.ReviewReq{Content: "", BookId: 1, UserId: 1}, err: review.ErrEmpty},
		{newReview: review.ReviewReq{Content: "Test", BookId: 1, UserId: 1}, err: nil},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := service.AddReview(t.Context(), tt.newReview)
			if !errors.Is(err, tt.err) {
				t.Fatalf("Error %v, correct %v", err, tt.err)
			}

			if tt.err != nil {
				t.Skip()
			}

			r, ok := repo.Reviews[repo.Id]
			if !ok {
				t.Fatalf("Review not saved")
			}

			if r.BookId != tt.newReview.BookId || r.UserId != tt.newReview.UserId || r.Content != tt.newReview.Content {
				t.Errorf("Error %v, correct %v", r, tt.newReview)
			}
		})
	}
}

func TestDeleteReview(t *testing.T) {
	repo := fakes.NewFakeReview()
	service := review.NewReviewService(repo)

	repo.Reviews[1] = review.Review{Id: 1, Content: "Test 1", UserId: 2, BookId: 3}
	repo.Reviews[2] = review.Review{Id: 2, Content: "Test 2", UserId: 4, BookId: 3}

	tests := []struct {
		reviewId, userId int
		err              error
	}{
		{reviewId: 1, userId: 6, err: review.ErrUnauth},
		{reviewId: 5, userId: 1, err: review.ErrNotFound},
		{reviewId: 1, userId: 2},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := service.DeleteReview(t.Context(), tt.reviewId, tt.userId)
			if !errors.Is(err, tt.err) {
				t.Fatalf("Error %v, correct %v", err, tt.err)
			}

			if tt.err != nil {
				t.Skip()
			}

			_, ok := repo.Reviews[tt.reviewId]
			if ok {
				t.Errorf("Review not deleted")
			}
		})
	}
}

func TestUpdateReview(t *testing.T) {
	repo := fakes.NewFakeReview()
	service := review.NewReviewService(repo)

	repo.Reviews[7] = review.Review{Id: 7, Content: "Test 7", UserId: 10, BookId: 3}
	repo.Reviews[8] = review.Review{Id: 8, Content: "Test 8", UserId: 14, BookId: 3}

	tests := []struct {
		newReview review.ReviewUpdate
		err       error
	}{
		{newReview: review.ReviewUpdate{}, err: review.ErrEmpty},
		{newReview: review.ReviewUpdate{Content: "..", UserId: 2, ReviewId: 1}, err: review.ErrNotFound},
		{newReview: review.ReviewUpdate{Content: "..", UserId: 2, ReviewId: 7}, err: review.ErrUnauth},
		{newReview: review.ReviewUpdate{Content: "..", UserId: 10, ReviewId: 7}, err: nil},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := service.UpdateReview(t.Context(), tt.newReview)
			if !errors.Is(err, tt.err) {
				t.Fatalf("Error %v, correct %v", err, tt.err)
			}

			if tt.err != nil {
				t.Skip()
			}

			r := repo.Reviews[tt.newReview.ReviewId]
			if r.Content != tt.newReview.Content {
				t.Errorf("Error %v, correct %v", r.Content, tt.newReview.Content)
			}
		})
	}
}

func TestReviewLike(t *testing.T) {
	repo := fakes.NewFakeReview()
	service := review.NewReviewService(repo)

	repo.Reviews[7] = review.Review{Id: 7, Content: "Test 7", UserId: 10, BookId: 3}

	tests := []struct {
		reviewId, userId int
		likeCount        int
		liked            bool
		err              error
	}{
		{reviewId: 1, err: review.ErrNotFound},
		{reviewId: 7, likeCount: 0, liked: false, err: nil},
		{reviewId: 7, likeCount: 1, liked: true, err: nil},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := service.ReviewLike(t.Context(), tt.reviewId, tt.userId, tt.liked)
			if !errors.Is(err, tt.err) {
				t.Fatalf("Error %v, correct %v", err, tt.err)
			}

			if tt.err != nil {
				t.Skip()
			}

			r := repo.Reviews[tt.reviewId]
			if r.Likes != tt.likeCount {
				t.Errorf("Error %v, correct %v", r.Likes, tt.likeCount)
			}
		})
	}
}

func TestReportReview(t *testing.T) {
	repo := fakes.NewFakeReview()
	service := review.NewReviewService(repo)

	repo.Reviews[7] = review.Review{Id: 7, Content: "Test 7", UserId: 10, BookId: 3}

	tests := []struct {
		reviewId int
		problem  string
		err      error
	}{
		{reviewId: 1, err: review.ErrNotFound},
		{reviewId: 7, problem: "Test", err: nil},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := service.ReportReview(t.Context(), tt.reviewId, tt.problem)
			if !errors.Is(err, tt.err) {
				t.Fatalf("Error %v, correct %v", err, tt.err)
			}

			if tt.err != nil {
				t.Skip()
			}

			if !slices.Contains(repo.Reports, tt.problem) {
				t.Errorf("Error %v, correct %v", repo.Reports, tt.problem)
			}
		})
	}
}
