package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/domain/port"
)

// Config holds the settings needed to construct the OpenAI chat adapter.
type Config struct {
	APIKey      string
	Model       string
	BaseURL     string
	Temperature float64
}

// openAILLM adapts a langchaingo OpenAI chat model to the domain port.LLM.
type openAILLM struct {
	model       *openai.LLM
	temperature float64
}

// compile-time check that openAILLM satisfies the port.
var _ port.LLM = (*openAILLM)(nil)

// NewOpenAI builds a port.LLM backed by langchaingo's OpenAI provider. The model
// name is configurable (e.g. gpt-4o-mini, gpt-4o, gpt-4-turbo).
func NewOpenAI(cfg Config) (port.LLM, error) {
	opts := []openai.Option{
		openai.WithToken(cfg.APIKey),
		openai.WithModel(cfg.Model),
	}
	if cfg.BaseURL != "" {
		opts = append(opts, openai.WithBaseURL(cfg.BaseURL))
	}

	model, err := openai.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("init openai llm: %w", err)
	}

	return &openAILLM{model: model, temperature: cfg.Temperature}, nil
}

func (l *openAILLM) Generate(ctx context.Context, prompt string) (string, error) {
	out, err := llms.GenerateFromSinglePrompt(ctx, l.model, prompt,
		llms.WithTemperature(l.temperature),
	)
	if err != nil {
		return "", fmt.Errorf("openai generate: %w", err)
	}
	return out, nil
}
