package main

import (
	"log"

	"github.com/gin-gonic/gin"
	docs "github.com/jared-scarr/portfolio-monorepo/packages/observability/docs"
	"github.com/jared-scarr/portfolio-monorepo/packages/observability/internal/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Metrics API
// @version 0.1
// @description Provides health checks and Prometheus metrics.
// @host localhost:8081
// @BasePath /
func main() {
	docs.SwaggerInfo.Title = "Metrics API"
	docs.SwaggerInfo.Description = "Provides health checks and Prometheus metrics."
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.BasePath = "/"
	r := gin.Default()
	r.Use(handlers.MetricsMiddleware())

	// Health
	r.GET("/health", handlers.Health)
	r.GET("/ready", handlers.Ready)

	// Metrics
	r.GET("/metrics", handlers.Metrics)

	// Swagger docs
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Metrics API running at http://localhost:8081")
	log.Println("Swagger docs available at http://localhost:8081/docs/index.html")

	r.Run(":8081")
}
