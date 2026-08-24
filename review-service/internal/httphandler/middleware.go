package httphandler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"review-service/internal/configs"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func JwtMiddleware(j *configs.JwtAuth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := r.Header.Get("Authorization")
			if t == "" {
				errorResponse(w, ErrNoToken)
				return
			}

			token := strings.TrimPrefix(t, "Bearer ")

			claims, err := j.CheckToken(token)
			if err != nil {
				errorResponse(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", claims.UserId)
			ctx = context.WithValue(ctx, "role", claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
