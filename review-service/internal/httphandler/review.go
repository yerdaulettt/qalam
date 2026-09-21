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

func (h *reviewHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.PathValue("username")

	reviews, err := h.service.GetUserReviews(r.Context(), username)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		errorResponse(w, err)
	}
}

func (h *reviewHandler) GetMyReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		errorResponse(w, ErrUserId)
		return
	}

	reviews, err := h.service.GetMyReviews(r.Context(), userId)
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

func (h *reviewHandler) ReviewLike(w http.ResponseWriter, r *http.Request) {
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

	liked, err := strconv.ParseBool(r.URL.Query().Get("liked"))
	if err != nil {
		liked = true
	}

	err = h.service.ReviewLike(r.Context(), reviewId, userId, liked)
	if err != nil {
		errorResponse(w, err)
		return
	}

	w.Write([]byte(`{"message": "ok"}`))
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
