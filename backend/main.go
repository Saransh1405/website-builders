package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"website-builder/internal/config"
	"website-builder/internal/handler"
	"website-builder/internal/repository"
	"website-builder/internal/router"
	"website-builder/internal/service"
)

func main() {
	// TODO: Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// TODO: Initialize template repository
	templateRepo, err := repository.NewFileTemplateRepository()
	if err != nil {
		log.Fatalf("Failed to initialize template repository: %v", err)
	}

	// TODO: Initialize Claude AI service
	aiService, err := service.NewClaudeAIService()
	if err != nil {
		log.Fatalf("Failed to initialize Claude AI service: %v", err)
	}

	// TODO: Initialize intent parser
	intentParser := service.NewIntentParserService()

	// TODO: Initialize post-processor
	postProcessor := service.NewPostProcessorService()

	// TODO: Initialize code generation service
	codeGenService := service.NewCodeGenerationService(
		templateRepo,
		aiService,
		intentParser,
		postProcessor,
	)

	// TODO: Initialize HTTP handler
	httpHandler := handler.NewHTTPHandler(codeGenService)

	// TODO: Setup router
	router := router.SetupRouter(cfg, httpHandler)

	// TODO: Create and start HTTP server
	server := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
