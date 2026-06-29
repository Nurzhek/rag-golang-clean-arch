package rest

import (
	"log/slog"
	"net/http"

	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/handler"
	"github.com/Nurzhek/rag-golang-clean-arch/internal/delivery/rest/middleware"
)

// NewRouter wires HTTP routes to handlers and applies global middleware. Route
// patterns use the method-aware mux available in Go 1.22+.
func NewRouter(docs *handler.DocumentHandler, query *handler.QueryHandler, jobs *handler.JobHandler, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.Health)

	// Documents
	mux.HandleFunc("POST /api/v1/documents", docs.Ingest)      // queue inline ingest -> 200 + job
	mux.HandleFunc("PUT /api/v1/documents", docs.IngestFile)   // chunked async file load -> 202 + job
	mux.HandleFunc("GET /api/v1/documents", docs.List)         // list documents
	mux.HandleFunc("DELETE /api/v1/documents/{id}", docs.Delete) // delete a document

	// Query
	mux.HandleFunc("POST /api/v1/query", query.Query)

	// Jobs (polling)
	mux.HandleFunc("GET /api/v1/jobs/{id}", jobs.Get)

	return middleware.Chain(mux,
		middleware.Recover(log),
		middleware.Logging(log),
	)
}
