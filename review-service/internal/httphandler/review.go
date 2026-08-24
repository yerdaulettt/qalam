package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"review-service/internal/review"
)

type reviewHandler struct {
	service *review.ReviewService
}

func newReviewHandler(s *review.ReviewService) *reviewHandler {
	return &reviewHandler{service: s}
}

func (h *reviewHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorResponse(w, ErrNumber)
		return
	}

	reviews, err := h.service.GetReviews(r.Context(), bookId)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		errorResponse(w, err)
	}
}

func (h *reviewHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorResponse(w, ErrNumber)
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		errorResponse(w, ErrUserId)
		return
	}

	var newReview review.ReviewReq
	err = json.NewDecoder(r.Body).Decode(&newReview)
	if err != nil {
		errorResponse(w, err)
		return
	}

	newReview.UserId = userId
	newReview.BookId = bookId

	rev, err := h.service.AddReview(r.Context(), newReview)
	if err != nil {
		errorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(rev)
	if err != nil {
		errorResponse(w, err)
	}
}

func (h *reviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reviewId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorResponse(w, ErrNumber)
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		errorResponse(w, ErrUserId)
		return
	}

	var newReview review.ReviewUpdate
	err = json.NewDecoder(r.Body).Decode(&newReview)
	if err != nil {
		errorResponse(w, err)
		return
	}

	newReview.ReviewId = reviewId
	newReview.UserId = userId

	rev, err := h.service.UpdateReview(r.Context(), newReview)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(rev)
	if err != nil {
		errorResponse(w, err)
	}
}

func (h *reviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reviewId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorResponse(w, ErrNumber)
		return
	}

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		errorResponse(w, ErrUserId)
		return
	}

	rev, err := h.service.DeleteReview(r.Context(), reviewId, userId)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(rev)
	if err != nil {
		errorResponse(w, err)
	}
}
