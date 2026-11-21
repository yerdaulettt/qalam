package httphandler

import (
	"net/http"

	"user-service/internal/auth"
	"user-service/internal/user"

	"github.com/go-chi/chi/v5"
)

func NewUserRouter(s *user.UserService, j *auth.JwtAuth) http.Handler {
	r := chi.NewRouter()

	userH := newUserHandler(s)

	r.Group(func(r chi.Router) {
		r.HandleFunc("POST /auth/register", userH.UserRegister)
		r.HandleFunc("POST /auth/login", userH.Login)
		r.HandleFunc("POST /auth/refresh", userH.TokenRefresh)
	})

	r.Group(func(r chi.Router) {
		r.Use(JwtMiddleware(j))
		r.HandleFunc("GET /my/profile", userH.GetMyProfile)
	})

	return r
}
