package auth

import (
	"fmt"
	"fpgschiba.com/automation-meal/database"
	"fpgschiba.com/automation-meal/router/roles"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"path/filepath"
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
			match, err := filepath.Match(route.Path, path)
			if err != nil {
				fmt.Println(err)
				return false
			}
			if match && slices.Contains(route.Methods, method) {
				return true
			}
		}
	}
	return false
}
