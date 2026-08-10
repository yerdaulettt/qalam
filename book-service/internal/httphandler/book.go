package httphandler

import (
	"encoding/json"
	"net/http"

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
