package users

import (
	"fpgschiba.com/automation-meal/database"
	"fpgschiba.com/automation-meal/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Register(c *gin.Context) {
	body := registerRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}
}

func AssignRole(c *gin.Context) {
	userId := c.Param("id")
	body := assignRoleRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}
	log.WithFields(log.Fields{
		"userId": userId,
	})
}

func ViewProfile(c *gin.Context) {
	userId := c.Param("id")
	log.WithFields(log.Fields{
		"userId": userId,
	})
}

func DeleteUser(c *gin.Context) {
	userId := c.Param("id")
	log.WithFields(log.Fields{
		"userId": userId,
	})
}

func UpdateUser(c *gin.Context) {
	userId := c.Param("id")
	body := updateUserRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}
	log.WithFields(log.Fields{
		"userId": userId,
	})
}

func ListUsers(c *gin.Context) {
	err, users := database.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponse(err))
		return
	}
	var userProfiles []UserProfile
	for _, user := range users {
		userProfiles = append(userProfiles, UserProfile{
			Id:                user.ID.Hex(),
			Email:             user.Email,
			DisplayName:       user.DisplayName,
			ProfilePictureUrl: user.ProfilePictureURL,
		})
	}
	c.JSON(http.StatusOK, ListUsersResponse{
		Response: util.Response{
			Status:  "success",
			Message: "Users retrieved successfully",
		},
		Users: userProfiles,
	})

}

func ResetPassword(c *gin.Context) {
	userId := c.Param("id")
	body := passwordResetRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}
	log.WithFields(log.Fields{
		"userId": userId,
	})
}
