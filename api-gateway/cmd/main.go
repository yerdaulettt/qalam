package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"api-gateway/internal/config"
	"api-gateway/internal/proxy"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	r, err := proxy.InitRoutes(c)
	if err != nil {
		log.Fatal(err)
	}

	notify, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	server := &http.Server{
		Addr:    ":8080",
		Handler: proxy.LogMiddleware(r),
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

	stopCtx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	defer stop()

	if err := server.Shutdown(stopCtx); err != nil {
		log.Fatal(err)
	}

	log.Println("Stopped")
}
