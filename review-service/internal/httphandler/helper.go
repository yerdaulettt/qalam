package httphandler

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrInternal = errors.New("Internal server error")
	ErrNoToken  = errors.New("No token")
)

func errorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	log.Println(err)

	switch err {
	case ErrInternal:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + ErrInternal.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"error":"` + err.Error() + `"}`))
}
