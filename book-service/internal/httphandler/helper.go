package httphandler

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrInternal = errors.New("Internal server error")
)

func errorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	log.Println(err)

	w.Write([]byte(`{"error":"` + ErrInternal.Error() + `"}`))
}
