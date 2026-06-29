package port

import "context"

// LLM is the port for large language model text generation. Implementations
// live in the infrastructure layer (e.g. an OpenAI-backed adapter).
type LLM interface {
	// Generate returns the model's completion for the given prompt.
	Generate(ctx context.Context, prompt string) (string, error)
}
