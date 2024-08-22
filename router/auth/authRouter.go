package auth

import (
	"crypto/sha256"
	"fmt"
	"fpgschiba.com/automation-meal/database"
	"fpgschiba.com/automation-meal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	body := loginRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}
	h := sha256.New()
	h.Write([]byte(body.Password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))
	err, exists, userId := database.DoesEmailExist(body.Email)
	if !exists || err != nil || userId == "" {
		c.JSON(http.StatusNotFound, util.GetResponseWithMessage(fmt.Sprintf("User with email: %s does not exist", body.Email)))
		return
	}
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetResponseWithMessage(fmt.Sprintf("User with id: %s did not have a correct ID.", userId)))
		return
	}
	err, passwordMatches, user := database.PasswordMatchesUser(userObjectId, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetResponseWithMessage(fmt.Sprintf("Password matching failed with follwing error: %s", err.Error())))
		return
	}
	if !passwordMatches || user == nil {
		c.JSON(http.StatusForbidden, util.GetResponseWithMessage(fmt.Sprintf("Incorrect Password for user: %s", body.Email)))
		return
	}
	claims := TokenClaims{
		Email: body.Email,
		ID:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}
	token, err := generateJWTToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetResponseWithMessage(fmt.Sprintf("Error generating token: %s", err.Error())))
		return
	}
	c.JSON(http.StatusOK, loginResponse{
		Response: util.Response{Status: "success", Message: "Login Successful."},
		Token:    token,
		User: User{
			Id:                userId,
			Email:             user.Email,
			DisplayName:       user.DisplayName,
			ProfilePictureUrl: user.ProfilePictureURL,
		},
	})
}

func ResetPassword(c *gin.Context) {
	body := resetPasswordRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, util.GetErrorResponse(err))
		return
	}
}
