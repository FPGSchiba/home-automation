package auth

import (
	"fmt"
	"fpgschiba.com/automation/database"
	"fpgschiba.com/automation/models"
	"fpgschiba.com/automation/util"
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
	err, exists, userId := database.DoesEmailExist(body.Email)
	if !exists || err != nil || userId == "" {
		c.JSON(http.StatusNotFound, util.GetErrorResponseWithMessage(fmt.Sprintf("User with email: %s does not exist", body.Email)))
		return
	}
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseWithMessage(fmt.Sprintf("User with id: %s did not have a correct ID.", userId)))
		return
	}
	err, passwordMatches, user := database.PasswordMatchesUser(userObjectId, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseWithMessage(fmt.Sprintf("Password matching failed with follwing error: %s", err.Error())))
		return
	}
	if !passwordMatches || user == nil {
		c.JSON(http.StatusForbidden, util.GetErrorResponseWithMessage(fmt.Sprintf("Incorrect Password for user: %s", body.Email)))
		return
	}
	config := util.Config{}
	config.GetConfig()
	claims := TokenClaims{
		Email: body.Email,
		ID:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Security.TokenExpiration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "schiba.com",
			Subject:   "Home Automation",
			ID:        userId,
		},
	}
	token, err := generateJWTToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseWithMessage(fmt.Sprintf("Error generating token: %s", err.Error())))
		return
	}
	err = database.InsertAuthEvent(models.EventLogin, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetErrorResponseWithMessage(fmt.Sprintf("Error logging Event: %s", err.Error())))
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
