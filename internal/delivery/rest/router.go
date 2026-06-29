package rest

import (
	"log/slog"
	"net/http"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/handler"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/middleware"
)

// NewRouter wires HTTP routes to handlers and applies global middleware. Route
// patterns use the method-aware mux available in Go 1.22+.
func NewRouter(rag *handler.RAGHandler, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /api/v1/documents", rag.Ingest)
	mux.HandleFunc("POST /api/v1/query", rag.Query)

	return middleware.Chain(mux,
		middleware.Recover(log),
		middleware.Logging(log),
	)
}
