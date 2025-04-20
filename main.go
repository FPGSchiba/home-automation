package main

import (
	"fpgschiba.com/automation/database"
	"fpgschiba.com/automation/router"
	"fpgschiba.com/automation/util/backup"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"strings"
)

func main() {
	engine := router.GetRouter()
	backup.StartScheduler()
	database.StartMigrations()

	if gin.Mode() == gin.ReleaseMode {
		engine.Use(static.Serve("/", static.LocalFile("./public/assets/", true)))
		engine.NoRoute(func(c *gin.Context) {
			if !strings.HasPrefix(c.Request.RequestURI, "/api") {
				c.File("./public/assets/index.html")
			}
			if !strings.HasPrefix(c.Request.RequestURI, "/auth") {
				c.File("./public/assets/index.html")
			}
		})
	}

	defer database.Disconnect()
	defer backup.StopScheduler()

	err := engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}
