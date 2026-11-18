package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ShortenURL(c *gin.Context) {
	linkStr, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	link, err := url.Parse(string(linkStr))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if link.Scheme != "http" && link.Scheme != "https" {
		c.Status(http.StatusBadRequest)
		return
	}

	if link.Hostname() == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	id, err := h.service.ShortenURL(string(linkStr))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	shortURL := fmt.Sprintf("http://%s/%s", c.Request.Host, id)
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusTemporaryRedirect, shortURL)
}
