package backups

import (
	"fpgschiba.com/automation/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBackupSettings(c *gin.Context) {
	// Get backup settings
	c.JSON(http.StatusNotImplemented, util.GetErrorResponseWithMessage("Not Implemented"))
}

func UpdateBackupSettings(c *gin.Context) {
	// Update backup settings
	c.JSON(http.StatusNotImplemented, util.GetErrorResponseWithMessage("Not Implemented"))
}
