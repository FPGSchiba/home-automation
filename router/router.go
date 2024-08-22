package router

import (
	"fpgschiba.com/automation-meal/router/auth"
	"fpgschiba.com/automation-meal/router/base"
	"fpgschiba.com/automation-meal/router/roles"
	"fpgschiba.com/automation-meal/router/users"
	"fpgschiba.com/automation-meal/util"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(util.JSONLogMiddleware())
	router.Use(util.CORS(util.CORSOptions{}))

	apiGroup := router.Group("/api/v1")
	{
		apiGroup.GET("/", base.GetVersion)
		apiGroup.Use(auth.AuthMiddleware())
		rolesGroup := apiGroup.Group("/roles")
		{
			rolesGroup.POST("/", roles.CreateRole)
		}
		usersGroup := apiGroup.Group("/users")
		{
			usersGroup.GET("/", users.ListUsers)
			usersGroup.POST("/register", users.Register)
			usersGroup.POST("/:id/role", users.AssignRole)
			usersGroup.GET("/:id", users.ViewProfile)
			usersGroup.DELETE("/:id", users.DeleteUser)
			usersGroup.PUT("/:id", users.UpdateUser)
			usersGroup.POST("/:id/reset-password", users.ResetPassword)
		}
	}
	router.POST("/auth/login", auth.Login)
	router.POST("/auth/reset-password", auth.ResetPassword)
	return router
}
