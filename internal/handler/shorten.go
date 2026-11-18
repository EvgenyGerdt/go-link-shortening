package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	linkStr, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, err := url.Parse(string(linkStr))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if link.Scheme != "http" && link.Scheme != "https" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if link.Hostname() == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.service.ShortenURL(string(linkStr))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := fmt.Sprintf("http://%s/%s", r.Host, id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(shortURL))
}
