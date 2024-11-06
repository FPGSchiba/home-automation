package permissions

import (
	"fpgschiba.com/automation-meal/router/roles"
	"fpgschiba.com/automation-meal/util"
)

type listPermissionsResponse struct {
	util.Response
	Permissions []roles.Permission `json:"permissions"`
}
