package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogoalish/atm/internal/logger"
	"go.uber.org/zap"
)

func RequestLogger(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		ctx := logger.WithLogger(c.Request.Context(), l)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		duration := time.Since(startTime)
		l.Info("Request details",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)
	}
}
