package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jared-scarr/portfolio-monorepo/apps/feature-flags-api/flags"
)

// -------- Types --------

type FlagStatus struct {
	Key     string `json:"key"`
	Enabled bool   `json:"enabled"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// -------- Handlers --------

// GetFlags godoc
// @Summary Get all flags for an environment
// @Produce json
// @Param   env  query   string  true  "Environment"  Enums(local,prod)
// @Success 200  {object}  map[string]bool
// @Failure 400  {object}  ErrorResponse
// @Failure 500  {object}  ErrorResponse
// @Router  /flags [get]
func GetFlags(c *gin.Context) {
	env := c.Query("env")
	if env != "local" && env != "prod" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid env; must be local or prod"})
		return
	}

	flagsMap, err := flags.GetAllFlags(env)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, flagsMap)
}

// GetFlagByKey godoc
// @Summary Get a single flag’s status for an environment
// @Produce json
// @Param   key  path    string  true  "Flag key (snake_case)"
// @Param   env  query   string  true  "Environment"  Enums(local,prod)
// @Success 200  {object}  FlagStatus
// @Failure 400  {object}  ErrorResponse
// @Failure 404  {object}  ErrorResponse
// @Failure 500  {object}  ErrorResponse
// @Router  /flags/{key} [get]
func GetFlagByKey(c *gin.Context) {
	env := c.Query("env")
	key := c.Param("key")

	if env != "local" && env != "prod" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid env; must be local or prod"})
		return
	}

	val, ok, err := flags.GetSingleFlag(env, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "unknown flag key"})
		return
	}

	c.JSON(http.StatusOK, FlagStatus{Key: key, Enabled: val})
}
