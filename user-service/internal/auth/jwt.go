package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	jwt.RegisteredClaims
	UserId int
	Role   string
	Type   string
}

type jwtAuth struct {
	secret     []byte
	accessTtl  time.Duration
	refreshTtl time.Duration
}

func NewJwtAuth(s []byte, aTtl, rTtl time.Duration) *jwtAuth {
	return &jwtAuth{
		secret:     s,
		accessTtl:  aTtl,
		refreshTtl: rTtl,
	}
}

func (j *jwtAuth) ParseToken(tokenString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrIncorrectToken
		}

		return j.secret, nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpired
		}

		return nil, ErrIncorrectToken
	}

	tokenClaims, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, ErrIncorrectToken
	}

	return tokenClaims, nil
}

func (j *jwtAuth) newAccessToken(userId int, role string) (string, error) {
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

func (j *jwtAuth) newRefreshToken(userId int, role string) (string, error) {
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

func (j *jwtAuth) newTokens(userId int, role string) (JwtTokens, error) {
	access, err := j.newAccessToken(userId, role)
	if err != nil {
		return JwtTokens{}, err
	}

	refresh, err := j.newRefreshToken(userId, role)
	if err != nil {
		return JwtTokens{}, err
	}

	return JwtTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (j *jwtAuth) refreshAccess(refresh string) (string, error) {
	claims, err := j.ParseToken(refresh)
	if err != nil {
		return "", err
	}

	if claims.Type != "refresh" {
		return "", ErrIncorrectToken
	}

	access, err := j.newAccessToken(claims.UserId, claims.Role)
	if err != nil {
		return "", ErrIncorrectToken
	}

	return access, nil
}
