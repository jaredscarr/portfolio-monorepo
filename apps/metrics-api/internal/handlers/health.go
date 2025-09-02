package handlers

import "github.com/gin-gonic/gin"

// Health godoc
// @Summary Liveness probe
// @Description Returns OK if service is alive
// @Tags health
// @Success 200 {object} map[string]string
// @Router /health [get]
func Health(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
}

// Ready godoc
// @Summary Readiness probe
// @Description Checks dependencies like Postgres/Redis
// @Tags health
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /ready [get]
func Ready(c *gin.Context) {
    // For now, always ready
    c.JSON(200, gin.H{"status": "ready"})
}
