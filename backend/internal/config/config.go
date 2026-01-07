package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Claude    ClaudeConfig
	CORS      CORSConfig
	Templates TemplatesConfig
	Logging   LoggingConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type ClaudeConfig struct {
	APIKey string
	Model  string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type TemplatesConfig struct {
	Dir string
}

type LoggingConfig struct {
	Level string
}

func LoadConfig() (*Config, error) {
	// TODO: Load environment variables from .env file
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Claude: ClaudeConfig{
			APIKey: getEnv("CLAUDE_API_KEY", ""),
			Model:  getEnv("CLAUDE_MODEL", "claude-3-5-sonnet-20241022"),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseStringSlice(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000")),
			AllowedMethods: parseStringSlice(getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS")),
			AllowedHeaders: parseStringSlice(getEnv("CORS_ALLOW_HEADERS", "Content-Type,Authorization")),
		},
		Templates: TemplatesConfig{
			Dir: getEnv("TEMPLATES_DIR", "./templates"),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	// TODO: Validate required configurations
	if config.Claude.APIKey == "" {
		return nil, fmt.Errorf("CLAUDE_API_KEY is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	for _, item := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
