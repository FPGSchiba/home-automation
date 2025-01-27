package backups

import (
	"fpgschiba.com/automation-meal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListBackups(c *gin.Context) {
	// List all backups
	c.JSON(http.StatusNotImplemented, util.GetResponseWithMessage("Not Implemented"))
}
