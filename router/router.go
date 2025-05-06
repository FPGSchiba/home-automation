package router

import (
	"fpgschiba.com/automation/router/auth"
	"fpgschiba.com/automation/router/backups"
	"fpgschiba.com/automation/router/base"
	"fpgschiba.com/automation/router/jobs"
	"fpgschiba.com/automation/router/permissions"
	"fpgschiba.com/automation/router/roles"
	"fpgschiba.com/automation/router/users"
	"fpgschiba.com/automation/util"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(util.JSONLogMiddleware())
	router.Use(util.CORS(util.CORSOptions{
		Origin: "*",
	}))

	apiGroup := router.Group("/api/v1")
	{
		apiGroup.GET("/", base.GetVersion)
		apiGroup.Use(auth.Middleware())
		rolesGroup := apiGroup.Group("/roles")
		{
			rolesGroup.POST("/", roles.CreateRole)
			rolesGroup.GET("/", roles.ListRoles)
			rolesGroup.PATCH("/:id", roles.UpdateRole)
			rolesGroup.DELETE("/:id", roles.DeleteRole)
			rolesGroup.GET("/:id", roles.GetRole)
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
		permissionsGroup := apiGroup.Group("/permissions")
		{
			permissionsGroup.GET("/", permissions.ListPermissions)
		}
		backupsGroup := apiGroup.Group("/backups")
		{
			backupsGroup.GET("/", backups.ListBackups)
			backupsGroup.GET("/:id/download", backups.DownloadBackup)
			backupsGroup.DELETE("/:id", backups.DeleteBackup)
			backupsGroup.GET("/:id", backups.GetBackupDetails)
			backupsGroup.GET("/:id/logs", backups.GetBackupLogs)
			backupsSettingsGroup := backupsGroup.Group("/settings")
			{
				backupsSettingsGroup.GET("/", backups.GetBackupSettings)
				backupsSettingsGroup.POST("/", backups.UpdateBackupSettings)
			}
		}
		apiGroup.GET("/job-types", jobs.GetJobTypes)
		jobsGroup := apiGroup.Group("/jobs")
		{
			jobsGroup.GET("/", jobs.ListJobs)
			jobsGroup.POST("/", jobs.CreateJob)
			jobsGroup.GET("/:id", jobs.GetJob)
			jobsGroup.PATCH("/:id", jobs.UpdateJob)
			jobsGroup.DELETE("/:id", jobs.DeleteJob)
		}
	}
	router.POST("/auth/login", auth.Login)
	router.POST("/auth/reset-password", auth.ResetPassword)

	// Configurations for the gin router
	router.RemoveExtraSlash = true
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return router
}
