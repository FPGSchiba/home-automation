package roles

import (
	"fpgschiba.com/automation/util"
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

func ListRoles(c *gin.Context) {
	return
}

func UpdateRole(c *gin.Context) {
	roleId := c.Param("id")
	body := roleUpdateRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}
	println(roleId)
}

func DeleteRole(c *gin.Context) {
	roleId := c.Param("id")
	println(roleId)
}

func GetRole(c *gin.Context) {
	roleId := c.Param("id")
	println(roleId)
}
