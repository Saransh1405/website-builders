package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware provides request logging
func LoggerMiddleware() gin.HandlerFunc {
	// TODO: Implement request logging middleware
	return gin.Logger()
}

// RecoveryMiddleware provides panic recovery
func RecoveryMiddleware() gin.HandlerFunc {
	// TODO: Implement panic recovery middleware
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v", recovered)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		c.Abort()
	})
}

// TraceIDMiddleware adds a trace ID to each request
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-ID")
		if traceID == "" {
			traceID = generateTraceID()
		}

		c.Set("traceId", traceID)
		c.Header("X-Request-ID", traceID)
		c.Next()
	}
}

// generateTraceID generates a simple trace ID
func generateTraceID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
