package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"review-service/internal/review"
)

type adminHandler struct {
	service *review.AdminService
}

func newAdminHandler(s *review.AdminService) *adminHandler {
	return &adminHandler{service: s}
}

func (h *adminHandler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reviews, err := h.service.GetAllReviews(r.Context())
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *adminHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reviewId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorResponse(w, errNumber)
		return
	}

	rev, err := h.service.DeleteReview(r.Context(), reviewId)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&rev)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *adminHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reports, err := h.service.GetReports(r.Context())
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&reports)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *adminHandler) DeleteFakeReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reviewId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorResponse(w, errNumber)
		return
	}

	err = h.service.DeleteFakeReports(r.Context(), reviewId)
	if err != nil {
		errorResponse(w, err)
		return
	}

	w.Write([]byte(`{"message": "ok"}`))
}
