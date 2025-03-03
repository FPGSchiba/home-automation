package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func getDurationInMilliseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}

func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := fmt.Sprintf("%.2fms", getDurationInMilliseconds(start))

		var fields log.Fields

		if c.FullPath() != "" {
			fields = log.Fields{
				"duration":    duration,
				"method":      c.Request.Method,
				"path":        c.Request.RequestURI,
				"status":      c.Writer.Status(),
				"referrer":    c.Request.Referer(),
				"api_version": ApiVersion,
			}
		} else {
			fields = log.Fields{
				"duration": duration,
				"method":   c.Request.Method,
				"path":     c.Request.RequestURI,
				"status":   c.Writer.Status(),
				"referrer": c.Request.Referer(),
			}
		}

		entry := log.WithFields(fields)

		if c.Writer.Status() >= 500 {
			entry.Error(c.Errors.String())
		} else if c.Writer.Status() >= 400 {
			entry.Warn("")
		} else {
			entry.Info("")
		}
	}
}
