package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {

	log.Printf("/healthcheck")
	message := fmt.Sprintf("Health check at %s", time.Now().Format(time.RFC3339))
	c.JSON(http.StatusOK, gin.H{"message": message})

}
