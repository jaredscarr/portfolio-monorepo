package main

import (
    "log"

    "github.com/gin-gonic/gin"
    _ "github.com/jared-scarr/portfolio-monorepo/apps/metrics-api/docs" // swagger docs
    "github.com/jared-scarr/portfolio-monorepo/apps/metrics-api/internal/handlers"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
)

// @title Metrics API
// @version 0.1
// @description Provides health checks and Prometheus metrics.
// @host localhost:8081
// @BasePath /
func main() {
    r := gin.Default()
	
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
