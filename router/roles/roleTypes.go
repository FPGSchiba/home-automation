package roles

type Permission struct {
	Routes []Route `json:"routes" binding:"required"`
}

type Route struct {
	Path        string                 `json:"path"`
	Methods     []string               `json:"methods"`
	JSONFilters map[string]interface{} `json:"jsonFilter"`
}

type roleCreationRequest struct {
	RoleName    string       `json:"roleName" binding:"required"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions" binding:"required"`
}

type roleUpdateRequest struct {
	RoleId      string       `json:"roleId" binding:"required"`
	RoleName    string       `json:"roleName"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}
