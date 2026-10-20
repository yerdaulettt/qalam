package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	p, err := mq.NewPublisher(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close(ctx)

	adminP := rabbitmq.NewAdminPublisher(p, mq.Management)

	jwtAuth, err := configs.NewJwtAuth()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(httphandler.LogMiddleware)

	bookR := postgres.NewBookRepo(db)
	bookS := book.NewBookService(bookR)
	bookRouter := httphandler.NewBookRouter(bookS, jwtAuth)

	adminR := postgres.NewAdminRepo(db)
	adminS := book.NewAdminService(adminR, adminP)
	adminRouter := httphandler.NewAdminRouter(adminS, jwtAuth)

	r.Mount("/api", bookRouter)
	r.Mount("/admin", adminRouter)

	server := &http.Server{
		Addr:    ":8082",
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
