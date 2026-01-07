# Website Builder Backend

A Golang backend service for AI-powered code generation with clean architecture.

## Project Structure

```
backend/
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Domain models and interfaces
│   ├── handler/         # HTTP handlers and middleware
│   ├── repository/      # Data access layer
│   ├── router/          # Route configuration
│   └── service/         # Business logic layer
├── templates/           # Template files for code generation
├── .gitignore
├── go.mod
├── main.go
└── README.md
```

## Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Create `.env` file:**
   ```env
   CLAUDE_API_KEY=your_claude_api_key_here
   PORT=8080
   GIN_MODE=debug
   ```

3. **Run the server:**
   ```bash
   go run main.go
   ```

## Implementation Guide

This is a starter project with placeholder implementations. You'll need to implement:

### Services (`internal/service/`)
- **ai_service.go**: Claude API integration with streaming
- **code_generation_service.go**: Main orchestration flow
- **intent_parser_service.go**: Parse user prompts to extract intent
- **post_processor_service.go**: Post-process generated code (add configs, etc.)

### Repository (`internal/repository/`)
- **template_repository.go**: Load template files from filesystem

### Handlers (`internal/handler/`)
- **http_handler.go**: HTTP request handlers for code generation endpoints
- **middleware.go**: Request logging, recovery, and trace ID middleware

### Router (`internal/router/`)
- **router.go**: Configure routes and middleware

## API Endpoints

Once implemented:

- `GET /health` - Health check
- `POST /api/v1/generate` - Generate code (non-streaming)
- `POST /api/v1/generate/stream` - Generate code with SSE streaming

## Next Steps

1. Implement each service layer starting with the domain interfaces
2. Add Claude API integration in `ai_service.go`
3. Implement intent parsing logic
4. Add template loading in repository
5. Complete HTTP handlers with SSE streaming
6. Test each component incrementally
