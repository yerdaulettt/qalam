package httphandler

import (
	"net/http"

	"book-service/internal/book"

	"github.com/go-chi/chi/v5"
)

func NewBookRouter(s *book.BookService) http.Handler {
	r := chi.NewRouter()

	bookH := NewBookHandler(s)

	r.HandleFunc("GET /", bookH.GetBooks)

	return r
}
