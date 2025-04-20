package permissions

import (
	"fpgschiba.com/automation/database"
	roles2 "fpgschiba.com/automation/router/roles"
	"fpgschiba.com/automation/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func ListPermissions(c *gin.Context) {
	userId := c.Keys["id"]
	userObjectId, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponseWithMessage("Invalid user ID"))
		return
	}
	err, roles := database.GetAllRolesForUser(userObjectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponse(err))
		return
	}

	allPermissions := make([]roles2.Permission, 0)
	for _, role := range roles {
		for _, permission := range role.Permissions {
			allPermissions = append(allPermissions, permission)
		}
	}
	c.JSON(http.StatusOK, listPermissionsResponse{
		Response: util.Response{
			Status:  "success",
			Message: "Users retrieved successfully",
		},
		Permissions: allPermissions,
	})
}
