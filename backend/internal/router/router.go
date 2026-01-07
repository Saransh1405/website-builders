package router

import (
	"website-builder/internal/config"
	"website-builder/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router
func SetupRouter(cfg *config.Config, httpHandler *handler.HTTPHandler) *gin.Engine {
	gin.SetMode(cfg.Server.GinMode)

	router := gin.New()

	// TODO: Add middleware
	// router.Use(handler.LoggerMiddleware())
	// router.Use(handler.RecoveryMiddleware())
	// router.Use(handler.TraceIDMiddleware())

	// TODO: Configure CORS
	// Uncomment and implement when ready:
	// import "github.com/gin-contrib/cors"
	// corsConfig := cors.Config{
	//     AllowOrigins:     cfg.CORS.AllowedOrigins,
	//     AllowMethods:     cfg.CORS.AllowedMethods,
	//     AllowHeaders:     cfg.CORS.AllowedHeaders,
	//     ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
	//     AllowCredentials: true,
	// }
	// router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", httpHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", httpHandler.HealthCheck)

		// TODO: Add code generation endpoints
		// v1.POST("/generate", httpHandler.GenerateCode)
		// v1.POST("/generate/stream", httpHandler.GenerateCodeStream)
	}

	return router
}
