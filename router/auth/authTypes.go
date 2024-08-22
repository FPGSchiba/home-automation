package auth

import (
	"fpgschiba.com/automation-meal/util"
	"github.com/golang-jwt/jwt/v5"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Id                string `json:"id"`
	Email             string `json:"email"`
	DisplayName       string `json:"displayName"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
}

type loginResponse struct {
	util.Response
	Token string `json:"token"`
	User  User   `json:"user"`
}

type resetPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type TokenClaims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}
