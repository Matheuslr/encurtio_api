package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheuslr/encurtio/internal/service"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{service: service}
}

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	url, err := h.service.Shorten(c.Request.Context(), req.URL)

	if err != nil {
		// attach error to context so middleware can log it
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"short_url": url.ShortCode,
	})
}
