package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	notify, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

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

	server := &http.Server{
		Addr:    ":8081",
		Handler: httphandler.LogMiddleware(r),
	}

	go func() {
		log.Println("Starting")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-notify.Done()
	cancel()
	log.Println("Shutdown signal")

	stopCtx, stop := context.WithTimeout(ctx, 10*time.Second)
	defer stop()

	if err := server.Shutdown(stopCtx); err != nil {
		log.Fatal(err)
	}

	log.Println("Stopped")
}
