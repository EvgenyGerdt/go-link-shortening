package handler

import (
	"fmt"
	"io"
	"net/http"
)

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.service.ShortenURL(string(url))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := fmt.Sprintf("http://%s/%s", r.Host, id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(shortURL))
}
