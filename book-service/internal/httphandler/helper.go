package httphandler

import (
	"errors"
	"log"
	"net/http"

	"book-service/internal/book"
)

var (
	ErrInternal = errors.New("Internal server error")
	errNumber   = errors.New("Incorrect number")
	errQuery    = errors.New("Incorrect query param")
	errUserId   = errors.New("Incorrect user id")
	errNoToken  = errors.New("No token")
	errRole     = errors.New("Incorrect role")
)

func errorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	log.Println(err)

	switch err {
	case book.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + ErrInternal.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"error":"` + err.Error() + `"}`))
}
