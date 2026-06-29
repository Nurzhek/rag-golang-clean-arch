// Command server is the composition root: it loads configuration, constructs the
// concrete adapters, wires them into the use cases, and serves the HTTP API.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Nurzhek/rag-golang-clean-arch/config"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/handler"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/infrastructure/embedding"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/infrastructure/llm"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/infrastructure/splitter"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/infrastructure/vectorstore"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/usecase"
	"github.com/Nurzhek/rag-golang-clean-arch/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(cfg.LogLevel)

	// --- Infrastructure: adapters implementing the domain ports. ---
	embedder, err := embedding.NewOpenAI(embedding.Config{
		APIKey:  cfg.OpenAIAPIKey,
		Model:   cfg.EmbeddingModel,
		BaseURL: cfg.OpenAIBaseURL,
	})
	if err != nil {
		return err
	}

	chat, err := llm.NewOpenAI(llm.Config{
		APIKey:      cfg.OpenAIAPIKey,
		Model:       cfg.LLMModel,
		BaseURL:     cfg.OpenAIBaseURL,
		Temperature: cfg.LLMTemperature,
	})
	if err != nil {
		return err
	}

	textSplitter := splitter.NewRecursive(splitter.Config{
		ChunkSize:    cfg.ChunkSize,
		ChunkOverlap: cfg.ChunkOverlap,
	})

	store := vectorstore.NewMemory()

	// --- Use cases: application business rules. ---
	ingestUC := usecase.NewIngestUseCase(textSplitter, embedder, store, newIDGenerator())
	queryUC := usecase.NewQueryUseCase(embedder, store, chat, usecase.DefaultPromptBuilder, cfg.RetrievalTopK)

	// --- Delivery: HTTP transport. ---
	ragHandler := handler.NewRAGHandler(ingestUC, queryUC)
	router := rest.NewRouter(ragHandler, log)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Info("server starting", "addr", srv.Addr, "llm_model", cfg.LLMModel, "embedding_model", cfg.EmbeddingModel)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Block until a termination signal or a fatal server error.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

// newIDGenerator returns a generator of random hex IDs for stored chunks.
func newIDGenerator() usecase.IDGenerator {
	return func() string {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return strconv.FormatInt(time.Now().UnixNano(), 16)
		}
		return hex.EncodeToString(b)
	}
}
