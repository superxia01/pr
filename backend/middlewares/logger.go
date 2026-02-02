package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%s] %s %s - Status: %d - Latency: %v",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			status,
			latency,
		)
	}
}
