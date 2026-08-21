package main

import (
	"context"
	"log"
	"net/http"

	"book-service/internal/book"
	"book-service/internal/configs"
	"book-service/internal/httphandler"
	"book-service/internal/postgres"
	"book-service/internal/rabbitmq"

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

	mq, err := rabbitmq.NewPublisherConn(ctx, "amqp://guest:guest@localhost:5673/")
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close(ctx)
	adminP := rabbitmq.NewAdminPublisher(mq)

	r := chi.NewRouter()

	r.Use(httphandler.LogMiddleware)

	bookR := postgres.NewBookRepo(db)
	bookS := book.NewBookService(bookR)
	bookRouter := httphandler.NewBookRouter(bookS)

	adminR := postgres.NewAdminRepo(db)
	adminS := book.NewAdminService(adminR, adminP)
	adminRouter := httphandler.NewAdminRouter(adminS)

	r.Mount("/api", bookRouter)
	r.Mount("/admin", adminRouter)

	log.Fatal(http.ListenAndServe(":8080", r))
}
