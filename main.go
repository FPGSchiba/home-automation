package main

import (
	"fpgschiba.com/automation-meal/database"
	"fpgschiba.com/automation-meal/router"
	"github.com/gin-gonic/contrib/static"
)

var frontendPaths = []string{
	"/",
	"/login",
	"/reset-password",
	"/not-found",
	"/about",
}

func main() {
	engine := router.GetRouter()

	database.StartMigrations()

	for _, path := range frontendPaths {
		engine.Use(static.Serve(path, static.LocalFile("./public/assets/", true)))
	}

	engine.Use(static.Serve("/login", static.LocalFile("./public/assets/", true)))

	defer database.Disconnect()

	err := engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}
