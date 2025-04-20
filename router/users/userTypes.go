package users

import "fpgschiba.com/automation/util"

type UserProfile struct {
	Id                string `json:"id"`
	Email             string `json:"email"`
	DisplayName       string `json:"displayName"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
}

type registerRequest struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
}

type assignRoleRequest struct {
	RoleId string `json:"role_id" binding:"required"`
}

type updateUserRequest struct {
	DisplayName       string `json:"displayName"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
}

type passwordResetRequest struct {
	NewPassword string `json:"newPassword" binding:"required"`
}

type ListUsersResponse struct {
	util.Response
	Users []UserProfile `json:"users"`
}
