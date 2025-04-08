package util

import "github.com/gin-gonic/gin"

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func GetResponse(message string, success bool) gin.H {
	status := "success"
	if !success {
		status = "error"
	}
	return gin.H{
		"message": message,
		"status":  status,
	}
}
