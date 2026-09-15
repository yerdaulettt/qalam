package httphandler

import (
	"net/http"

	"user-service/internal/auth"

	"github.com/go-chi/chi/v5"
)

func NewAuthRouter(s *auth.AuthService) http.Handler {
	r := http.NewServeMux()

	authH := NewAuthHandler(s)

	r.HandleFunc("POST /register", authH.UserRegister)
	r.HandleFunc("POST /login", authH.Login)
	r.HandleFunc("POST /refresh", authH.TokenRefresh)

	return r
}

func NewProfileRouter(s *auth.AuthService, j *auth.JwtAuth) http.Handler {
	r := chi.NewRouter()

	r.Use(JwtMiddleware(j))

	authH := NewAuthHandler(s)
	r.HandleFunc("GET /my/profile", authH.GetMyProfile)

	return r
}
