package main

import (
	"context"
	"log"
	"net/http"

	"review-service/internal/configs"
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

	postgresCfg, err := configs.NewPostgresConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewPostgresPool(ctx, postgresCfg.NewDBUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mq, err := rabbitmq.NewConsumerConn(ctx, "amqp://guest:guest@localhost:5673/")
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close(ctx)

	reviewR := postgres.NewReviewRepo(db)
	reviewS := review.NewReviewService(reviewR)

	rConsumer := rabbitmq.NewReviewConsumer(mq, reviewS)
	go rConsumer.DeleteBookId(ctx)

	r := chi.NewRouter()

	log.Fatal(http.ListenAndServe(":8082", r))
}
