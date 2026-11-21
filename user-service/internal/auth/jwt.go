package auth

import (
	"errors"
	"time"

	"user-service/internal/user"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	jwt.RegisteredClaims
	UserId int
	Role   string
	Type   string
}

type JwtAuth struct {
	secret     []byte
	accessTtl  time.Duration
	refreshTtl time.Duration
}

func NewJwtAuth(s []byte, aTtl, rTtl time.Duration) *JwtAuth {
	return &JwtAuth{
		secret:     s,
		accessTtl:  aTtl,
		refreshTtl: rTtl,
	}
}

func (j *JwtAuth) ParseToken(tokenString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, user.ErrIncorrectToken
		}

		return j.secret, nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, user.ErrExpired
		}

		return nil, user.ErrIncorrectToken
	}

	tokenClaims, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, user.ErrIncorrectToken
	}

	return tokenClaims, nil
}

func (j *JwtAuth) newAccessToken(userId int, role string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTtl)),
		},
		UserId: userId,
		Role:   role,
		Type:   "access",
	})

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtAuth) newRefreshToken(userId int, role string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTtl)),
		},
		UserId: userId,
		Role:   role,
		Type:   "refresh",
	})

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtAuth) NewTokens(userId int, role string) (user.JwtTokens, error) {
	access, err := j.newAccessToken(userId, role)
	if err != nil {
		return user.JwtTokens{}, err
	}

	refresh, err := j.newRefreshToken(userId, role)
	if err != nil {
		return user.JwtTokens{}, err
	}

	return user.JwtTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (j *JwtAuth) RefreshAccess(refresh string) (string, error) {
	claims, err := j.ParseToken(refresh)
	if err != nil {
		return "", err
	}

	if claims.Type != "refresh" {
		return "", user.ErrIncorrectToken
	}

	access, err := j.newAccessToken(claims.UserId, claims.Role)
	if err != nil {
		return "", user.ErrIncorrectToken
	}

	return access, nil
}
