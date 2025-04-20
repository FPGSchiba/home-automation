package jobs

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestMock(t *testing.T) {
	r := SetUpRouter()
	r.GET("/", ListJobs)

	assert.Equal(t, 1, 1)
}
