package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration, sourced from environment variables.
type Config struct {
	HTTPPort string
	LogLevel string

	// OpenAI credentials are shared by both the chat model and the embedder.
	OpenAIAPIKey  string
	OpenAIBaseURL string

	LLMModel       string
	LLMTemperature float64
	EmbeddingModel string

	ChunkSize     int
	ChunkOverlap  int
	RetrievalTopK int
}

// Load reads configuration from the environment, applying defaults and
// validating required values. A local .env file is loaded if present.
func Load() (*Config, error) {
	_ = godotenv.Load() // optional: ignore when no .env file exists

	cfg := &Config{
		HTTPPort:       getEnv("HTTP_PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		OpenAIAPIKey:   os.Getenv("OPENAI_API_KEY"),
		OpenAIBaseURL:  os.Getenv("OPENAI_BASE_URL"),
		LLMModel:       getEnv("LLM_MODEL", "gpt-4o-mini"),
		LLMTemperature: getEnvFloat("LLM_TEMPERATURE", 0.2),
		EmbeddingModel: getEnv("EMBEDDING_MODEL", "text-embedding-3-small"),
		ChunkSize:      getEnvInt("CHUNK_SIZE", 1000),
		ChunkOverlap:   getEnvInt("CHUNK_OVERLAP", 200),
		RetrievalTopK:  getEnvInt("RETRIEVAL_TOP_K", 4),
	}

	if cfg.OpenAIAPIKey == "" {
		return nil, errors.New("OPENAI_API_KEY is required")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}
