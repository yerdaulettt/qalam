package httphandler

import (
	"errors"
	"log"
	"net/http"

	"user-service/internal/auth"
)

var (
	ErrJson = errors.New("Incorrect json")
)

func errorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	log.Println(err)

	switch err {
	case auth.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case ErrJson, auth.ErrEmptyFields, auth.ErrShortPassword, auth.ErrIncorrectPassword, auth.ErrIncorrectToken:
		w.WriteHeader(http.StatusBadRequest)
	case auth.ErrUsername:
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal error"}`))
		return
	}

	w.Write([]byte(`{"error":"` + err.Error() + `"}`))
}
