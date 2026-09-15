package httphandler

import (
	"errors"
	"log"
	"net/http"

	"user-service/internal/auth"
)

var (
	errJson    = errors.New("Incorrect json")
	errUserId  = errors.New("Incorrect user id")
	errRole    = errors.New("Incorrect role")
	errNoToken = errors.New("No token")
	errToken   = errors.New("Incorrect token")
)

func errorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	log.Println(err)

	switch err {
	case auth.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case errJson, errToken, errUserId, auth.ErrEmptyFields, auth.ErrShortPassword, auth.ErrIncorrectPassword, auth.ErrIncorrectToken:
		w.WriteHeader(http.StatusBadRequest)
	case auth.ErrUsername:
		w.WriteHeader(http.StatusConflict)
	case errRole:
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal error"}`))
		return
	}

	w.Write([]byte(`{"error":"` + err.Error() + `"}`))
}
