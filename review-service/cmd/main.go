package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"review-service/internal/configs"
	"review-service/internal/httphandler"
	"review-service/internal/postgres"
	"review-service/internal/rabbitmq"
	"review-service/internal/review"

	"github.com/go-chi/chi/v5"
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

	db, err := postgres.NewPostgresPool(ctx, postgresCfg.NewDBUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mq, err := rabbitmq.NewConn(ctx, configs.NewRabbitUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close(ctx)

	reviewR := postgres.NewReviewRepo(db)
	reviewS := review.NewReviewService(reviewR)

	adminR := postgres.NewAdminRepo(db)
	adminS := review.NewAdminService(adminR)

	jwtAuth, err := configs.NewJwtAuth()
	if err != nil {
		log.Fatal(err)
	}

	consume, err := mq.NewConsumer(ctx, "book.*", "books")
	if err != nil {
		log.Fatal(err)
	}
	rConsumer := rabbitmq.NewReviewConsumer(consume, reviewS)
	go rConsumer.DeleteBookId(ctx)

	r := chi.NewRouter()
	r.Use(httphandler.LogMiddleware)

	r.Mount("/api", httphandler.NewReviewRouter(reviewS, jwtAuth))
	r.Mount("/admin", httphandler.NewAdminRouter(adminS, jwtAuth))

	server := &http.Server{
		Addr:    ":8083",
		Handler: r,
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
