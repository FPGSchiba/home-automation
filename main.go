package main

import (
	"fpgschiba.com/automation-meal/router"
	"github.com/gin-gonic/contrib/static"
)

func main() {
	engine := router.GetRouter()

	engine.Use(static.Serve("/", static.LocalFile("./public/assets/", true)))

	err := engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}
