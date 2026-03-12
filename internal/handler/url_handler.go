package handler

import (
	"fmt"
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
		message := fmt.Sprintln("failed to shorten url: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": message,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"short_url": url.ShortCode,
	})
}
