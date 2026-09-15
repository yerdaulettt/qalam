package main

import (
	"context"
	"log"
	"net/http"

	"user-service/internal/auth"
	"user-service/internal/configs"
	"user-service/internal/httphandler"
	"user-service/internal/postgres"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	postgresCfg, err := configs.NewPostgresConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewPostgresPool(ctx, postgresCfg)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	r := http.NewServeMux()

	jwtCfg, err := configs.NewJwtConfig()
	if err != nil {
		log.Fatal(jwtCfg)
	}

	authR := postgres.NewAuthRepo(db)
	jwtAuth := auth.NewJwtAuth(jwtCfg.Secret, jwtCfg.AccessTtl, jwtCfg.RefreshTtl)
	authS := auth.NewAuthService(authR, jwtAuth)
	authRouter := httphandler.NewAuthRouter(authS)

	profileRouter := httphandler.NewProfileRouter(authS, jwtAuth)

	r.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	r.Handle("/api/", http.StripPrefix("/api", profileRouter))

	log.Fatal(http.ListenAndServe(":8080", httphandler.LogMiddleware(r)))
}
