package roles

type Permission struct {
	Paths       []string               `json:"paths" binding:"required"`
	JSONFilters map[string]interface{} `json:"jsonFilters"`
}

type roleCreationRequest struct {
	RoleName    string       `json:"roleName" binding:"required"`
	Permissions []Permission `json:"permissions" binding:"required"`
}
