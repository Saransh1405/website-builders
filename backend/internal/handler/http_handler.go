package handler

import (
	"net/http"
	"time"
	"website-builder/internal/domain"
	"website-builder/internal/service"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	codeGenService *service.CodeGenerationService
}

func NewHTTPHandler(codeGenService *service.CodeGenerationService) *HTTPHandler {
	return &HTTPHandler{
		codeGenService: codeGenService,
	}
}

func (h *HTTPHandler) GenerateCodeStream(c *gin.Context) {
	// TODO: Implement SSE streaming handler
	// 1. Parse and validate request body
	// 2. Set up SSE headers
	// 3. Call codeGenService.GenerateCodeStream
	// 4. Stream responses back to client

	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *HTTPHandler) GenerateCode(c *gin.Context) {
	// TODO: Implement non-streaming handler
	// 1. Parse and validate request body
	// 2. Call codeGenService.GenerateCode
	// 3. Return JSON response

	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *HTTPHandler) HealthCheck(c *gin.Context) {
	status := domain.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Services: map[string]string{
			"api":       "operational",
			"claude":    "operational",
			"templates": "operational",
		},
	}

	c.JSON(http.StatusOK, status)
}

type ErrorDetail struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
	Value string `json:"value,omitempty"`
}

func (h *HTTPHandler) respondWithError(
	c *gin.Context,
	statusCode int,
	code string,
	message string,
	details []ErrorDetail,
) {
	traceID := c.GetString("traceId")
	if traceID == "" {
		traceID = c.GetString("X-Request-ID")
	}

	errorResponse := ErrorResponse{
		Error: ErrorInfo{
			Code:      code,
			Message:   message,
			Details:   details,
			TraceID:   traceID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	c.JSON(statusCode, errorResponse)
}

type ErrorResponse struct {
	Error ErrorInfo `json:"error"`
}

type ErrorInfo struct {
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Details   []ErrorDetail `json:"details,omitempty"`
	TraceID   string        `json:"traceId,omitempty"`
	Timestamp string        `json:"timestamp"`
}
