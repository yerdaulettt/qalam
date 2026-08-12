package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"book-service/internal/book"
)

type bookHandler struct {
	service *book.BookService
}

func NewBookHandler(s *book.BookService) *bookHandler {
	return &bookHandler{service: s}
}

func (h *bookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := h.service.GetBooks(r.Context())
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&books)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *bookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorResponse(w, err)
		return
	}

	b, err := h.service.GetBook(r.Context(), bookId)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(b)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *bookHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	g, err := h.service.GetGenres(r.Context())
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&g)
	if err != nil {
		errorResponse(w, err)
		return
	}
}
