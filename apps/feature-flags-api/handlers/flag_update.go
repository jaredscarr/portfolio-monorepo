package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jared-scarr/portfolio-monorepo/apps/feature-flags-api/flags"
)

// UpdateFlagRequest represents the request to update a flag
type UpdateFlagRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateFlag godoc
// @Summary Update a feature flag value dynamically
// @Produce json
// @Param   key  path    string  true  "Flag key (snake_case)"
// @Param   env  query   string  true  "Environment"  Enums(local,prod)
// @Param   request body UpdateFlagRequest true "Flag update request"
// @Success 200  {object}  FlagStatus
// @Failure 400  {object}  ErrorResponse
// @Failure 404  {object}  ErrorResponse
// @Failure 500  {object}  ErrorResponse
// @Router /admin/flags/{key} [put]
func UpdateFlag(c *gin.Context) {
	env := c.Query("env")
	key := c.Param("key")

	if env != "local" && env != "prod" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid env; must be local or prod"})
		return
	}

	var req UpdateFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Debug logging
	fmt.Printf("DEBUG: Updating flag %s in env %s to enabled=%v\n", key, env, req.Enabled)

	// Check if flag exists first
	_, exists, err := flags.GetSingleFlag(env, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "flag not found"})
		return
	}

	// Update the flag in memory
	err = flags.UpdateFlag(env, key, req.Enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, FlagStatus{Key: key, Enabled: req.Enabled})
}
