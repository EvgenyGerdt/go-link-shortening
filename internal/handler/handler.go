package handler

import "github.com/EvgenyGerdt/go-link-shortening.git/internal/service"

type URLService interface {
	ResolveURL(string) (string, error)
	ShortenURL(string) (string, error)
}

type Handler struct {
	service URLService
}

func NewHandler(s *service.Shortener) *Handler {
	return &Handler{
		service: s,
	}
}
