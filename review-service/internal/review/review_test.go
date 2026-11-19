package review

import (
	"database/sql"
	"errors"
	"testing"
)

func TestAddReview(t *testing.T) {
	repo := newFakeReview()
	service := NewReviewService(repo)

	tests := []struct {
		newReview ReviewReq
		repoErr   error
		err       error
	}{
		{newReview: ReviewReq{Content: ""}, err: ErrEmpty},
		{newReview: ReviewReq{Content: "Test"}, err: ErrEmpty},
		{newReview: ReviewReq{Content: "Test", BookId: 1}, err: ErrEmpty},
		{newReview: ReviewReq{Content: "Test", UserId: 1}, err: ErrEmpty},
		{newReview: ReviewReq{Content: "", BookId: 1, UserId: 1}, err: ErrEmpty},
		{newReview: ReviewReq{Content: "Test", BookId: 1, UserId: 1}, err: nil},
		{newReview: ReviewReq{Content: "Test", BookId: 1, UserId: 1}, repoErr: sql.ErrConnDone, err: sql.ErrConnDone},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			repo.err[0] = tt.repoErr

			err := service.AddReview(t.Context(), tt.newReview)
			if !errors.Is(err, tt.err) {
				t.Fatalf("Error %v, correct %v", err, tt.err)
			}

			if tt.err != nil {
				t.Skip()
			}

			r, ok := repo.reviews[repo.id]
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
	repo := newFakeReview()
	service := NewReviewService(repo)

	repo.reviews[1] = Review{Id: 1, Content: "Test 1", UserId: 2, BookId: 3}
	repo.reviews[2] = Review{Id: 2, Content: "Test 2", UserId: 4, BookId: 3}

	tests := []struct {
		reviewId, userId int
		repoErr          error
		userIdErr        error
		err              error
	}{
		{reviewId: 1, userId: 6, err: ErrUnauth},
		{reviewId: 5, userId: 1, err: ErrNotFound},
		{reviewId: 2, userId: 1, userIdErr: sql.ErrConnDone, err: sql.ErrConnDone},
		{reviewId: 1, userId: 2, repoErr: sql.ErrConnDone, err: sql.ErrConnDone},
		{reviewId: 1, userId: 2},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			repo.err[0] = tt.repoErr
			repo.err[1] = tt.userIdErr

			err := service.DeleteReview(t.Context(), tt.reviewId, tt.userId)
			if !errors.Is(err, tt.err) {
				t.Fatalf("Error %v, correct %v", err, tt.err)
			}

			if tt.err != nil {
				t.Skip()
			}

			_, ok := repo.reviews[tt.reviewId]
			if ok {
				t.Errorf("Review not deleted")
			}
		})
	}
}
