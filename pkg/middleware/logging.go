package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/HoangQuan74/goodie-api/pkg/logger"
	"go.uber.org/zap"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		requestID, _ := c.Get("request_id")

		l := logger.Get()
		l.Info("http request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.Any("request_id", requestID),
		)
	}
}
