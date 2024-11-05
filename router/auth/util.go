package auth

import (
	"fpgschiba.com/automation-meal/database"
	"fpgschiba.com/automation-meal/router/roles"
	"github.com/bmatcuk/doublestar"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"slices"
)

func verifyPermissions(userId string, path string, method string) bool {
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return false
	}
	err, userRoles := database.GetAllRolesForUser(userObjectId)
	if err != nil {
		return false
	}
	for _, role := range userRoles {
		if hasPermission(role.Permissions, path, method) {
			return true
		}
	}
	return false
}

func hasPermission(permissions []roles.Permission, path string, method string) bool {
	for _, permission := range permissions {
		for _, route := range permission.Routes {
			match, err := doublestar.Match(route.Path, path) // needed support for ** to match separators as well
			println(match)
			if err != nil {
				return false
			}
			if match && slices.Contains(route.Methods, method) {
				return true
			}
		}
	}
	return false
}
