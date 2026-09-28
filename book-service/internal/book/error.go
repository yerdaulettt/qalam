package book

import "errors"

var (
	ErrNotFound       = errors.New("Not found")
	ErrGenreNotFound  = errors.New("Genre not found")
	ErrAuthorNotFound = errors.New("Author not found")
)
