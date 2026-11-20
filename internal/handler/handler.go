package handler

type URLService interface {
	ResolveURL(string) (string, error)
	ShortenURL(string) (string, error)
}

type Handler struct {
	service URLService
	baseURL string
}

func NewHandler(s URLService, baseURL string) *Handler {
	return &Handler{
		service: s,
		baseURL: baseURL,
	}
}
