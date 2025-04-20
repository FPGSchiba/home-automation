package permissions

import (
	"fpgschiba.com/automation/router/roles"
	"fpgschiba.com/automation/util"
)

type listPermissionsResponse struct {
	util.Response
	Permissions []roles.Permission `json:"permissions"`
}
