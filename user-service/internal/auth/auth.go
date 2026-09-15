package auth

import (
	"errors"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserVerify struct {
	Id   int
	Role string
	Hash string
}

type JwtTokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type RegisterReq struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var (
	ErrEmptyFields       = errors.New("All fields required")
	ErrIncorrectToken    = errors.New("Incorrect token")
	ErrExpired           = errors.New("Token expired")
	ErrNotFound          = errors.New("Not found")
	ErrUsername          = errors.New("Username exists")
	ErrShortPassword     = errors.New("Minimum password len is 8")
	ErrIncorrectPassword = errors.New("Incorrect password")
)
