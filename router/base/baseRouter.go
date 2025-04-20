package base

import (
	"fpgschiba.com/automation/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": util.ApiVersion,
	})
}
