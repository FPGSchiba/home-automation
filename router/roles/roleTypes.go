package roles

type Permission struct {
	Routes []Route `json:"paths" binding:"required"`
}

type Route struct {
	Path        string                 `json:"path"`
	Methods     []string               `json:"methods"`
	JSONFilters map[string]interface{} `json:"jsonFilter"`
}

type roleCreationRequest struct {
	RoleName    string       `json:"roleName" binding:"required"`
	Permissions []Permission `json:"permissions" binding:"required"`
}

type roleUpdateRequest struct {
	RoleId      string       `json:"roleId" binding:"required"`
	RoleName    string       `json:"roleName"`
	Permissions []Permission `json:"permissions"`
}
