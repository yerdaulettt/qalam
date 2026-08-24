package configs

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	jwt.RegisteredClaims
	UserId int
	Role   string
	Type   string
}

type JwtAuth struct {
	secret []byte
}

func NewJwtAuth() (*JwtAuth, error) {
	s := []byte(os.Getenv("JWT_SECRET"))
	if len(s) == 0 {
		return nil, ErrEnv
	}

	return &JwtAuth{secret: s}, nil
}

func (j *JwtAuth) CheckToken(t string) (*claims, error) {
	token, err := jwt.ParseWithClaims(t, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrToken
		}

		return j.secret, nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenTime
		}

		return nil, ErrToken
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, ErrToken
	}

	if c.Type != "access" {
		return nil, ErrToken
	}

	return c, nil
}
