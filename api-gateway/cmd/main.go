package main

import (
	"log"
	"net/http"

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

	log.Fatal(http.ListenAndServe(":8080", proxy.LogMiddleware(r)))
}
