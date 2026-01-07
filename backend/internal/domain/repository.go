package domain

import (
	"context"
)

type TemplateRepository interface {
	LoadTemplate(ctx context.Context, templatePath string) ([]byte, error)
	ListTemplates(ctx context.Context) ([]string, error)
	TemplateExists(ctx context.Context, templatePath string) bool
}

type AIService interface {
	GenerateCodeStream(ctx context.Context, prompt string, systemPrompt string, templateContent string) (<-chan string, <-chan error)
	GenerateCode(ctx context.Context, prompt string, systemPrompt string, templateContent string) (string, error)
}

type IntentParser interface {
	ParseIntent(ctx context.Context, prompt string) (*Intent, error)
}

type PostProcessor interface {
	ProcessCode(ctx context.Context, code string, intent *Intent, options RequestOptions) (*GenerationResponse, error)
}
