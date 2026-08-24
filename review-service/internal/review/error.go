package review

import "errors"

var (
	ErrNotFound = errors.New("Not found")
	ErrEmpty    = errors.New("Empty field")
	ErrUnauth   = errors.New("Unauthorized")
)
