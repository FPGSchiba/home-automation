package roles

import (
	"fpgschiba.com/automation-meal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRole(c *gin.Context) {
	body := roleCreationRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}
}
