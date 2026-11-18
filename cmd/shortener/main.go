package main

import (
	"net/http"

	"github.com/EvgenyGerdt/go-link-shortening.git/internal/handler"
	"github.com/EvgenyGerdt/go-link-shortening.git/internal/repository"
	"github.com/EvgenyGerdt/go-link-shortening.git/internal/service"
)

func main() {
	mux := http.NewServeMux()

	// Repositories
	linkRepo := repository.NewLinkRepository()
	// Services
	shortenerService := service.NewShortener(linkRepo)
	// Handlers
	h := handler.NewHandler(shortenerService)

	mux.HandleFunc("POST /", h.ShortenURL)
	mux.HandleFunc("GET /{id}", h.ResolveURL)
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
