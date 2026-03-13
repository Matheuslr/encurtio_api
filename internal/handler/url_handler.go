package handler

import (
	"log"
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
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"short_url": h.service.BuildShortURL(url.ShortCode),
	})
}

func (h *URLHandler) GetRedirectURL(c *gin.Context) {
	code := c.Param("code")
	log.Printf("Received request for code: %s", code)
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing short code"})
		return
	}

	url, err := h.service.GetOriginalURL(c.Request.Context(), code)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "short code not found"})
		return
	}

	c.Redirect(http.StatusFound, url.Original)
}
