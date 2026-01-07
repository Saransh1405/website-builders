package domain

import "time"

type GenerationRequest struct {
	Prompt        string         `json:"prompt" binding:"required"`
	ComponentType string         `json:"componentType,omitempty"` // react, vue, etc.
	Options       RequestOptions `json:"options,omitempty"`
}

type RequestOptions struct {
	UseTypeScript bool   `json:"useTypeScript,omitempty"`
	IncludeTests  bool   `json:"includeTests,omitempty"`
	StyleLibrary  string `json:"styleLibrary,omitempty"` // tailwind, css, styled-components
}

type GenerationResponse struct {
	ComponentCode string            `json:"componentCode"`
	PackageJSON   string            `json:"packageJson,omitempty"`
	ConfigFiles   map[string]string `json:"configFiles,omitempty"`
	Files         []File            `json:"files,omitempty"`
}

type File struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Type    string `json:"type"` // component, config, test, etc.
}

type GenerationStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`   // pending, processing, completed, failed
	Progress  int       `json:"progress"` // 0-100
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type Intent struct {
	Type          string   `json:"type"` // react-component, vue-component, general, etc.
	ComponentName string   `json:"componentName,omitempty"`
	Features      []string `json:"features,omitempty"` // hooks, state, props, etc.
	Confidence    float64  `json:"confidence"`         // 0.0 - 1.0
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services,omitempty"`
}
