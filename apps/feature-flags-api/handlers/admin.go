package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jared-scarr/portfolio-monorepo/apps/feature-flags-api/flags"
)

// ReloadFlags godoc
// @Summary Reload flags from disk (internal use)
// @Produce json
// @Success 200  {object}  map[string]string
// @Failure 500  {object}  ErrorResponse
// @Router /admin/reload [post]
func ReloadFlags(c *gin.Context) {
	errLocal := flags.LoadFlagsFromDisk("local")
	errProd := flags.LoadFlagsFromDisk("prod")

	if errLocal != nil || errProd != nil {
		msg := "failed to reload flags:"
		if errLocal != nil {
			msg += " local=" + errLocal.Error()
		}
		if errProd != nil {
			msg += " prod=" + errProd.Error()
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "flags reloaded"})
}
