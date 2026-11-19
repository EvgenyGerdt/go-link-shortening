package main

import (
	"log"

	"github.com/EvgenyGerdt/go-link-shortening.git/internal/config"
	"github.com/EvgenyGerdt/go-link-shortening.git/internal/handler"
	"github.com/EvgenyGerdt/go-link-shortening.git/internal/repository"
	"github.com/EvgenyGerdt/go-link-shortening.git/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()
	router := gin.Default()

	// Repositories
	linkRepo := repository.NewLinkRepository()
	// Services
	shortenerService := service.NewShortener(linkRepo)
	// Handlers
	h := handler.NewHandler(shortenerService, cfg.BaseURL)

	router.POST("/", h.ShortenURL)
	router.GET("/:id", h.ResolveURL)
	router.GET("/", h.ResolveURL)

	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatal(err)
	}
}
