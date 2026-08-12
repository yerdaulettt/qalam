package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"book-service/internal/book"
)

type adminHandler struct {
	service *book.AdminService
}

func NewAdminHandler(s *book.AdminService) *adminHandler {
	return &adminHandler{service: s}
}

func (h *adminHandler) NewGenre(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var g struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		errorResponse(w, err)
		return
	}

	genre, err := h.service.NewGenre(r.Context(), g.Name)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&genre)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *adminHandler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	genreId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorResponse(w, err)
		return
	}

	var g struct {
		Name string `json:"name"`
	}

	err = json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		errorResponse(w, err)
		return
	}

	genre, err := h.service.UpdateGenre(r.Context(), genreId, g.Name)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&genre)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *adminHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	genreId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorResponse(w, err)
		return
	}

	genre, err := h.service.DeleteGenre(r.Context(), genreId)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&genre)
	if err != nil {
		errorResponse(w, err)
		return
	}
}
