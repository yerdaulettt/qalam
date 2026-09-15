package httphandler

import (
	"errors"
	"log"
	"net/http"

	"review-service/internal/review"
)

var (
	ErrInternal = errors.New("Internal server error")
	ErrNumber   = errors.New("Incorrect number")
	ErrUserId   = errors.New("Incorrect user id")
	ErrNoToken  = errors.New("No token")
	errRole     = errors.New("Incorrect role")
)

func errorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	log.Println(err)

	switch err {
	case review.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case review.ErrEmpty:
		w.WriteHeader(http.StatusBadRequest)
	case review.ErrUnauth:
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + ErrInternal.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"error":"` + err.Error() + `"}`))
}
