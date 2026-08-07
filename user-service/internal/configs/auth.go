package configs

import (
	"os"
	"strconv"
	"time"
)

type jwtConfig struct {
	Secret     []byte
	AccessTtl  time.Duration
	RefreshTtl time.Duration
}

func NewJwtConfig() (*jwtConfig, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return nil, ErrEnv
	}

	aTtl, err := strconv.Atoi(os.Getenv("JWT_ACCESS_TTL"))
	if err != nil {
		aTtl = 15
	}

	rTtl, err := strconv.Atoi(os.Getenv("JWT_REFRESH_TTL"))
	if err != nil {
		rTtl = 24
	}

	cfg := jwtConfig{
		Secret:     secret,
		AccessTtl:  time.Minute * time.Duration(aTtl),
		RefreshTtl: time.Hour * time.Duration(rTtl),
	}

	return &cfg, nil
}
