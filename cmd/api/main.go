package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/matheuslr/encurtio/configs"
	"github.com/matheuslr/encurtio/internal/database"
	"github.com/matheuslr/encurtio/internal/handler"
	"github.com/matheuslr/encurtio/internal/middleware"
	"github.com/matheuslr/encurtio/internal/repository"
	"github.com/matheuslr/encurtio/internal/service"
)

func main() {
	cfg := configs.Load()

	session, err := database.NewCassandraSession(database.CassandraConfig{
		Hosts:    cfg.Cassandra.Hosts,
		Keyspace: cfg.Cassandra.Keyspace,
	})

	if err != nil {
		log.Fatal("failed to connect to cassandra: ", err)
	}
	defer session.Close()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorCapture())
	//repos
	repo := repository.NewCassandraURLRepository(session)

	//services
	service := service.NewURLService(repo, *cfg)

	//handlres
	healthHandler := handler.NewHealthHandler()
	urlHandler := handler.NewURLHandler(service)

	//routes
	router.POST("/api/v1/url/shorten", urlHandler.Shorten)
	router.GET("/:code", urlHandler.GetRedirectURL)
	router.GET("/api/v1/health", healthHandler.Health)

	router.Run(":" + cfg.API.Port)
}
