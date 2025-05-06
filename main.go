package main

import (
	"flag"
	"fpgschiba.com/automation/database"
	"fpgschiba.com/automation/router"
	"fpgschiba.com/automation/util"
	"fpgschiba.com/automation/util/backup"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
)

func main() {
	engine := router.GetRouter()

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	var configFilepath string
	flag.StringVar(&configFilepath, "config", "config.yaml", "The Path to the config file")
	flag.Parse()

	util.SetConfigFilePath(configFilepath)

	backup.StartScheduler()
	database.StartMigrations()

	// TODO: See if we keep this maybe rather use nginx container
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

	if gin.Mode() == gin.ReleaseMode {
		// Run on port 80 for production
		conf := util.Config{}
		conf.GetConfig()

		if conf.TLS.Enabled {
			panic("TLS is not supported yet")
		} else {
			err := engine.Run(":80")
			if err != nil {
				panic(err)
			}
		}
	}
	err := engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}
