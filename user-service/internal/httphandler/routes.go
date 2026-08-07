package httphandler

import (
	"net/http"

	"user-service/internal/auth"
)

func NewAuthRouter(s *auth.AuthService) http.Handler {
	r := http.NewServeMux()

	authH := NewAuthHandler(s)

	r.HandleFunc("POST /register", authH.UserRegister)
	r.HandleFunc("POST /login", authH.Login)
	r.HandleFunc("POST /refresh", authH.TokenRefresh)

	return r
}
