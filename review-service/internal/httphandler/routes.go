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

	r.Group(func(r chi.Router) {
		r.HandleFunc("GET /books/{id}/reviews", reviewH.GetReviews)
		r.HandleFunc("GET /users/{username}/reviews", reviewH.GetUserReviews)
		r.HandleFunc("POST /reviews/{id}/reports", reviewH.ReportReview)
	})

	r.Group(func(r chi.Router) {
		r.Use(JwtMiddleware(jwtAuth))
		r.HandleFunc("GET /my/reviews", reviewH.GetMyReviews)
		r.HandleFunc("POST /books/{id}/reviews", reviewH.AddReview)
		r.HandleFunc("PUT /reviews/{id}/likes", reviewH.ReviewLike)
		r.HandleFunc("PATCH /reviews/{id}", reviewH.UpdateReview)
		r.HandleFunc("DELETE /reviews/{id}", reviewH.DeleteReview)
	})

	return r
}

func NewAdminRouter(service *review.AdminService, jwtAuth *configs.JwtAuth) http.Handler {
	r := chi.NewRouter()

	r.Use(JwtMiddleware(jwtAuth))
	r.Use(RoleMiddleware("admin"))

	adminH := newAdminHandler(service)
	r.HandleFunc("GET /reviews", adminH.GetAllReviews)
	r.HandleFunc("DELETE /reviews/{id}", adminH.DeleteReview)
	r.HandleFunc("GET /reports", adminH.GetReports)
	r.HandleFunc("DELETE /reviews/{id}/reports", adminH.DeleteFakeReports)

	return r
}
