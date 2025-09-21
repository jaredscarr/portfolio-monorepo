// @title       Feature Flags API
// @version     1.0.0
// @description Read-only boolean feature flags served from repo JSON (local/prod).
// @BasePath    /
// @schemes     http
// @produce     json
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/jared-scarr/portfolio-monorepo/apps/feature-flags-api/docs"
	"github.com/jared-scarr/portfolio-monorepo/apps/feature-flags-api/flags"
	"github.com/jared-scarr/portfolio-monorepo/apps/feature-flags-api/handlers"
	observability "github.com/jared-scarr/portfolio-monorepo/packages/observability/handlers"
)

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
