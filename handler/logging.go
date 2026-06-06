package handler

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

const requestIDHeader = "X-Request-ID"

// RequestID attaches a request id to every request, reusing an incoming
// X-Request-ID header when the caller already provides one. The id is echoed
// back on the response and made available to the logger.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = newRequestID()
		}
		c.Set("request_id", id)
		c.Writer.Header().Set(requestIDHeader, id)
		c.Next()
	}
}

// RequestLogger emits one structured log line per request once it completes.
// Any errors handlers attach via c.Error are logged with the same request id so
// failures behind a generic client response are still diagnosable.
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		requestID := c.GetString("request_id")

		if len(c.Errors) > 0 {
			logger.Error("request failed",
				"request_id", requestID,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"errors", c.Errors.String(),
			)
		}

		logger.Info("request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

func newRequestID() string {
	return randomHex(16)
}
