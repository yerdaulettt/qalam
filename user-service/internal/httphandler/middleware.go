package httphandler

import (
	"context"
	"log"
	"net/http"
	"slices"
	"strings"

	"user-service/internal/auth"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func JwtMiddleware(j *auth.JwtAuth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := r.Header.Get("Authorization")
			if t == "" {
				errorResponse(w, errNoToken)
				return
			}

			token := strings.TrimPrefix(t, "Bearer ")

			claims, err := j.ParseToken(token)
			if err != nil {
				errorResponse(w, err)
				return
			}

			if claims.Type != "access" {
				errorResponse(w, errToken)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", claims.UserId)
			ctx = context.WithValue(ctx, "role", claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RoleMiddleware(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value("role").(string)
			if !ok {
				errorResponse(w, errRole)
				return
			}

			if !slices.Contains(roles, role) {
				errorResponse(w, errRole)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
