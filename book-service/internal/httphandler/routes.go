package httphandler

import (
	"net/http"

	"book-service/internal/book"

	"github.com/go-chi/chi/v5"
)

func NewBookRouter(s *book.BookService) http.Handler {
	r := chi.NewRouter()

	bookH := NewBookHandler(s)

	r.HandleFunc("GET /books", bookH.GetBooks)
	r.HandleFunc("GET /books/{id}", bookH.GetBook)
	r.HandleFunc("GET /genres", bookH.GetGenres)

	return r
}

func NewAdminRouter(s *book.AdminService) http.Handler {
	r := chi.NewRouter()

	adminH := NewAdminHandler(s)

	r.HandleFunc("POST /genres", adminH.NewGenre)
	r.HandleFunc("PATCH /genres/{id}", adminH.UpdateGenre)
	r.HandleFunc("DELETE /genres/{id}", adminH.DeleteGenre)

	return r
}
