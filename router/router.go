package router

import (
	"example.com/automation-meal/router/base"
	"example.com/automation-meal/util"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(util.JSONLogMiddleware())
	router.Use(util.CORS(util.CORSOptions{}))

	superGroup := router.Group("/api/v1")
	{
		superGroup.GET("/", base.GetVersion)
	}
	return router
}
