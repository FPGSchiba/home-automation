package auth

import (
	"errors"
	"fmt"
	"fpgschiba.com/automation-meal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.GetResponseWithMessage("Authorization header is empty"))
			return
		}

		authToken := strings.Replace(authHeader, "Bearer ", "", 1)
		fmt.Println(authToken)
		if claims, err := VerifyJWTToken(authToken); err != nil {
			log.WithFields(log.Fields{
				"authentication": "true",
				"component":      "AuthMiddleware",
				"method":         c.Request.Method,
				"path":           c.Request.URL.Path,
				"error":          err.Error(),
			}).Error()
			if errors.Is(err, jwt.ErrECDSAVerification) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, util.GetResponseWithMessage("Token is invalid"))
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.GetResponseWithMessage(fmt.Sprintf("Failed token Validation: %s", err.Error())))
			return

		} else {
			fmt.Println(claims.Email)
			c.Set("email", claims.Email)
			c.Set("id", claims.ID)
			// TODO: Permission checking

			c.Next()
		}
	}
}
