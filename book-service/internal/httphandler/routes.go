package httphandler

import (
	"net/http"

	"book-service/internal/book"
	"book-service/internal/configs"

	"github.com/go-chi/chi/v5"
)

func NewBookRouter(s *book.BookService, jwtAuth *configs.JwtAuth) http.Handler {
	r := chi.NewRouter()

	bookH := NewBookHandler(s)

	r.Group(func(r chi.Router) {
		r.HandleFunc("GET /books", bookH.GetBooks)
		r.HandleFunc("GET /books/{id}", bookH.GetBook)
		r.HandleFunc("GET /genres", bookH.GetGenres)
	})

	r.Group(func(r chi.Router) {
		r.Use(JwtMiddleware(jwtAuth))
		r.HandleFunc("PUT /books/{id}/ratings", bookH.AddRating)
	})

	return r
}

func NewAdminRouter(s *book.AdminService, jwtAuth *configs.JwtAuth) http.Handler {
	r := chi.NewRouter()

	adminH := NewAdminHandler(s)
	r.Use(JwtMiddleware(jwtAuth))
	r.Use(RoleMiddleware("admin"))

	r.HandleFunc("POST /genres", adminH.NewGenre)
	r.HandleFunc("PATCH /genres/{id}", adminH.UpdateGenre)
	r.HandleFunc("DELETE /genres/{id}", adminH.DeleteGenre)
	r.HandleFunc("POST /books", adminH.AddBook)
	r.HandleFunc("PATCH /books/{id}", adminH.UpdateBook)
	r.HandleFunc("DELETE /books/{id}", adminH.DeleteBook)
	r.HandleFunc("POST /authors", adminH.AddAuthor)
	r.HandleFunc("PATCH /authors/{id}", adminH.UpdateAuthor)

	return r
}
