package main

import (
	"context"
	"log"
	"net/http"

	"book-service/internal/book"
	"book-service/internal/configs"
	"book-service/internal/httphandler"
	"book-service/internal/postgres"

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

	r := chi.NewRouter()

	r.Use(httphandler.LogMiddleware)

	bookR := postgres.NewBookRepo(db)
	bookS := book.NewBookService(bookR)
	bookRouter := httphandler.NewBookRouter(bookS)

	r.Mount("/api/books", bookRouter)

	log.Fatal(http.ListenAndServe(":8080", r))
}
