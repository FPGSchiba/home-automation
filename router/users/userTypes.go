package users

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
