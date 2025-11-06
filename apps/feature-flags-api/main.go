package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/jared-scarr/portfolio-monorepo/apps/feature-flags-api/docs"
	"github.com/jared-scarr/portfolio-monorepo/apps/feature-flags-api/flags"
	"github.com/jared-scarr/portfolio-monorepo/apps/feature-flags-api/handlers"
	observability "github.com/jared-scarr/portfolio-monorepo/packages/observability/handlers"
)

// @title       Feature Flags API
// @version     1.0.0
// @description Read-only boolean feature flags served from repo JSON (local/prod).
// @BasePath    /
// @schemes     http
// @produce     json

// getCORSOrigins returns the list of allowed CORS origins from environment variable or defaults
func getCORSOrigins() []string {
	defaultOrigins := []string{"http://localhost:3000", "http://portfolio:3000"}

	if corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); corsOrigins != "" {
		// Parse comma-separated origins
		origins := strings.Split(corsOrigins, ",")
		result := make([]string, 0, len(origins))
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}

	return defaultOrigins
}

func main() {
	if err := flags.LoadFlagsFromDisk("local"); err != nil {
		log.Fatal(err)
	}
	if err := flags.LoadFlagsFromDisk("prod"); err != nil {
		log.Fatal(err)
	}

	flags.LoadFlagsFromDisk("local")
	flags.LoadFlagsFromDisk("prod")

	docs.SwaggerInfo.Title = "Feature Flags API"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = "/"

	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     getCORSOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add observability middleware
	r.Use(observability.MetricsMiddleware())

	// Observability endpoints
	r.GET("/health", observability.Health)
	r.GET("/ready", observability.Ready)
	r.GET("/metrics", observability.Metrics)

	// Feature flags endpoints
	r.GET("/flags", handlers.GetFlags)
	r.GET("/flags/:key", handlers.GetFlagByKey)
	r.POST("/admin/reload", handlers.ReloadFlags)
	r.PUT("/admin/flags/:key", handlers.UpdateFlag)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("feature-flags-api listening on :4000")
	if err := r.Run(":4000"); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
