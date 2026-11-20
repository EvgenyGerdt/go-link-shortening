package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ResolveURL(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	url, err := h.service.ResolveURL(id)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Header("Location", url)
	c.Status(http.StatusTemporaryRedirect)
}
