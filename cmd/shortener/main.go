package main

import (
	"github.com/EvgenyGerdt/go-link-shortening.git/internal/handler"
	"github.com/EvgenyGerdt/go-link-shortening.git/internal/repository"
	"github.com/EvgenyGerdt/go-link-shortening.git/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Repositories
	linkRepo := repository.NewLinkRepository()
	// Services
	shortenerService := service.NewShortener(linkRepo)
	// Handlers
	h := handler.NewHandler(shortenerService)

	router.POST("/", h.ShortenURL)
	router.GET("/{id}", h.ResolveURL)
	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
