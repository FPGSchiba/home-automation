package backups

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"testing"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestMock(t *testing.T) {
	r := SetUpRouter()
	r.GET("/", ListBackups)

	assert.Equal(t, 1, 1)
}
