package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/config"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/gates"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/handlers"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/storage"
	observability "github.com/jared-scarr/portfolio-monorepo/packages/observability/handlers"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := storage.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize storage layer
	store := storage.NewOutboxStore(db)

	// Initialize feature flag client and simulation gates
	flagsClient := gates.NewHTTPFeatureFlagClient("http://localhost:4000")
	simulationGates := gates.NewSimulationGates(flagsClient, "local")

	// Initialize handlers
	h := handlers.New(store, cfg, simulationGates)

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://portfolio:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add observability middleware
	router.Use(observability.MetricsMiddleware())

	// Health check endpoints
	router.GET("/health", observability.Health)
	router.GET("/ready", observability.Ready)
	router.GET("/metrics", observability.Metrics)

	// API endpoints
	api := router.Group("/api/v1")
	{
		api.POST("/events", h.CreateEvent)
		api.GET("/events", h.ListEvents)
		api.GET("/events/:id", h.GetEvent)
		api.POST("/events/:id/retry", h.RetryEvent)
		api.DELETE("/events/:id", h.DeleteEvent)
	}

	// Admin endpoints
	admin := router.Group("/admin")
	{
		admin.POST("/publish", h.PublishEvents)
		admin.GET("/stats", h.GetStats)
		admin.GET("/simulation-status", h.GetSimulationStatus)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting outbox-api server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
