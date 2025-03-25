package backups

import (
	"fpgschiba.com/automation-meal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// This is useful for storing files: https://www.mongodb.com/docs/drivers/go/v1.8/fundamentals/gridfs/

func ListBackups(c *gin.Context) {
	// List all backups
	c.JSON(http.StatusNotImplemented, util.GetResponseWithMessage("Not Implemented"))
}

func DeleteBackup(c *gin.Context) {
	// Delete a backup
	c.JSON(http.StatusNotImplemented, util.GetResponseWithMessage("Not Implemented"))
}

func DownloadBackup(c *gin.Context) {
	// Download a backup
	c.JSON(http.StatusNotImplemented, util.GetResponseWithMessage("Not Implemented"))
}
