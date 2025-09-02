package handlers

import "github.com/gin-gonic/gin"

// Metrics godoc
// @Summary Prometheus metrics
// @Description Exposes metrics in Prometheus text format
// @Tags metrics
// @Produce plain
// @Success 200 {string} string "Prometheus metrics text format"
// @Router /metrics [get]
func Metrics(c *gin.Context) {
    c.String(200, "# HELP http_requests_total Total HTTP requests\nhttp_requests_total 42")
}
