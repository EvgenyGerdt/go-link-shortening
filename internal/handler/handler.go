package handler

import "github.com/EvgenyGerdt/go-link-shortening.git/internal/service"

type Handler struct {
	service *service.Shortener
}

func NewHandler(s *service.Shortener) *Handler {
	return &Handler{
		service: s,
	}
}
