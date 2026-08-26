package httphandler

import (
	"net/http"

	"review-service/internal/configs"
	"review-service/internal/review"

	"github.com/go-chi/chi/v5"
)

func NewReviewRouter(service *review.ReviewService, jwtAuth *configs.JwtAuth) http.Handler {
	r := chi.NewRouter()

	reviewH := newReviewHandler(service)

	r.Use(JwtMiddleware(jwtAuth))

	r.HandleFunc("GET /my/reviews", reviewH.GetMyReviews)
	r.HandleFunc("GET /users/{username}/reviews", reviewH.GetUserReviews)
	r.HandleFunc("GET /books/{id}/reviews", reviewH.GetReviews)
	r.HandleFunc("POST /books/{id}/reviews", reviewH.AddReview)
	r.HandleFunc("PATCH /reviews/{id}", reviewH.UpdateReview)
	r.HandleFunc("DELETE /reviews/{id}", reviewH.DeleteReview)

	return r
}
