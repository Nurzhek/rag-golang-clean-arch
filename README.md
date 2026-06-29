# RAG Backend (Go · Clean Architecture · langchaingo)

A production-style **Retrieval-Augmented Generation (RAG)** backend written in Go.
It exposes an HTTP API to **ingest documents** (synchronously queued or via chunked
async file upload), **manage** them, and **ask grounded questions**: documents are
chunked, embedded, and stored as vectors; questions retrieve the most relevant chunks
and feed them to an LLM that answers **only** from that context.

Built with [**langchaingo**](https://github.com/tmc/langchaingo) using the **OpenAI**
provider for both chat completion and embeddings, and structured around
**Clean Architecture** with **SOLID**, **DRY**, and **KISS** as first-class constraints.

---

## Features

- 🧩 **Clean Architecture** - domain, application, infrastructure, and delivery layers with a strict inward dependency rule.
- 🔌 **Pluggable by design** - LLM, embedder, vector store, text splitter, and job repository are interfaces (ports); swap implementations without touching business logic.
- 🤖 **OpenAI via langchaingo** - any OpenAI chat model (configurable through `LLM_MODEL`) plus `text-embedding-3-small` embeddings.
- ⚡ **Async ingestion with polling** - `POST` returns immediately with a job ID; `PUT` accepts a raw file and ingests it in **batched/chunked background** processing; poll job progress via `GET /jobs/{id}`.
- 🗂️ **Document management** - list ingested documents and delete a document (and all its chunks) by ID.
- 🗃️ **Zero-setup vector store** - dependency-free in-memory cosine-similarity store; perfect for local dev and tests, trivially replaceable with pgvector/Qdrant.
- 🌐 **Standard-library HTTP** - method-aware routing (Go 1.22+ `ServeMux`), structured logging (`log/slog`), panic recovery, and graceful shutdown.
- 🧪 **Tested core** - use cases and the store are unit-tested with fakes (no network required).
- 🐳 **Container-ready** - multi-stage `Dockerfile` (distroless) and `docker-compose.yml`.

---

## Architecture

The codebase follows **Clean Architecture**: dependencies point **inward only**.
Inner layers know nothing about outer layers — the domain has zero imports from
infrastructure or transport.

```
┌──────────────────────────────────────────────────────────────────┐
│  Delivery (HTTP)            cmd/server, internal/delivery/rest     │  ← frameworks & I/O
│   router · document/query/job handlers · middleware · DTOs         │
│        │ depends on                                                │
│        ▼                                                           │
│  Application (Use Cases)    internal/usecase                       │  ← business workflows
│   IngestUseCase · DocumentUseCase · JobUseCase · QueryUseCase      │
│        │ depends on                                                │
│        ▼                                                           │
│  Domain (Core)              internal/domain                        │  ← entities & ports
│   entities · ports (LLM, Embedder, VectorStore,                    │
│                     TextSplitter, JobRepository)                   │
│        ▲ implemented by                                            │
│        │                                                           │
│  Infrastructure (Adapters)  internal/infrastructure               │  ← langchaingo, stores
│   OpenAI LLM · OpenAI Embedder · Splitter · in-memory vector store │
│   · in-memory job store                                            │
└──────────────────────────────────────────────────────────────────┘
```

**The Dependency Rule.** `domain` imports nothing from the project. `usecase` imports
only `domain`. `infrastructure` implements `domain` ports. `delivery` depends on
`usecase` interfaces. The **composition root** (`cmd/server/main.go`) is the only place
that knows every concrete type — it wires them together.

### How the principles map to the code

| Principle | Where you see it |
|-----------|------------------|
| **S**ingle Responsibility | One concern per type: `asyncIngestInteractor` runs the ingest pipeline, `Memory` stores vectors, `DocumentHandler`/`QueryHandler`/`JobHandler` each translate HTTP for one resource. |
| **O**pen/Closed | Add a pgvector store, a Redis job store, or a different LLM by implementing a port — no existing code changes. `var _ port.X = (*T)(nil)` compile-time checks guard the contracts. |
| **L**iskov Substitution | Any `port.VectorStore`/`port.JobRepository` is interchangeable; the use cases are oblivious. |
| **I**nterface Segregation | Small, focused ports (`LLM`, `Embedder`, `VectorStore`, `TextSplitter`, `JobRepository`) instead of one fat "AI service" interface. |
| **D**ependency Inversion | Use cases and handlers depend on interfaces; concretes are injected at the composition root. |
| **DRY** | Shared JSON/error helpers and a single domain-error->HTTP mapping, one ingest pipeline reused by `POST` and `PUT`, one prompt builder, one config loader, reusable middleware chain. |
| **KISS** | Standard-library HTTP and routing, an in-memory store and job repo, env-based config — no framework ceremony. |

---

## Project structure

```
rag-golang-clean-arch/
├── cmd/
│   └── server/
│       └── main.go                  # composition root: wire & serve
├── config/
│   └── config.go                    # env-based configuration + validation
├── internal/
│   ├── domain/                      # enterprise core (no outward imports)
│   │   ├── entity/                  # Document, ScoredDocument, DocumentSummary, Job
│   │   ├── port/                    # LLM, Embedder, VectorStore, TextSplitter, JobRepository
│   │   └── errors.go                # sentinel domain errors
│   ├── usecase/                     # application business rules
│   │   ├── ingest.go                # async submit + background pipeline (split → embed → store)
│   │   ├── documents.go             # list / delete documents
│   │   ├── jobs.go                  # job status (polling read side)
│   │   ├── query.go                 # embed → retrieve → prompt → generate
│   │   └── prompt.go                # grounded RAG prompt builder
│   ├── infrastructure/              # adapters implementing the ports
│   │   ├── llm/openai.go            # langchaingo OpenAI chat
│   │   ├── embedding/openai.go      # langchaingo OpenAI embeddings
│   │   ├── splitter/recursive.go    # langchaingo recursive splitter
│   │   ├── vectorstore/memory.go    # in-memory cosine store
│   │   └── jobstore/memory.go       # in-memory job repository
│   └── delivery/
│       └── rest/                    # HTTP transport
│           ├── router.go            # routes + middleware
│           ├── handler/             # document, query, job, health, shared responders
│           ├── middleware/          # logging, panic recovery
│           └── dto/                 # request/response shapes
├── pkg/
│   └── logger/                      # slog setup
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

---

## Tech stack

- **Go 1.23** (uses `log/slog` and the method-aware `http.ServeMux`)
- **langchaingo** - LLM, embeddings, and text-splitting abstractions
- **OpenAI** - chat completion (`gpt-4o-mini` by default) and `text-embedding-3-small`
- **godotenv** - `.env` loading for local development

---

## Getting started

### Prerequisites

- Go **1.23+** ([install](https://go.dev/dl/))
- An **OpenAI API key**

### 1. Clone & configure

```bash
git clone https://github.com/Nurzhek/rag-golang-clean-arch.git
cd rag-golang-clean-arch

cp .env.example .env
# edit .env and set OPENAI_API_KEY
```

### 2. Install dependencies

```bash
make tidy            # or: go mod tidy
```

### 3. Run

```bash
make run             # or: go run ./cmd/server
```

The server listens on `http://localhost:8080`.

---

## Configuration

All configuration is read from environment variables (a local `.env` is loaded automatically).

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | `8080` | HTTP listen port |
| `LOG_LEVEL` | `info` | `debug` \| `info` \| `warn` \| `error` |
| `OPENAI_API_KEY` | - | **Required.** OpenAI API key |
| `OPENAI_BASE_URL` | _(unset)_ | Override base URL (Azure OpenAI / proxy / compatible server) |
| `LLM_MODEL` | `gpt-4o-mini` | Any OpenAI chat model (`gpt-4o`, `gpt-4-turbo`, `gpt-3.5-turbo`, …) |
| `LLM_TEMPERATURE` | `0.2` | Sampling temperature for generation |
| `EMBEDDING_MODEL` | `text-embedding-3-small` | OpenAI embedding model |
| `CHUNK_SIZE` | `1000` | Characters per chunk |
| `CHUNK_OVERLAP` | `200` | Overlap between adjacent chunks |
| `RETRIEVAL_TOP_K` | `4` | Number of chunks retrieved per query |
| `EMBED_BATCH_SIZE` | `64` | Chunks embedded+stored per batch during async ingestion |

---

## API reference

| Method & path | Purpose |
|---|---|
| `GET /health` | Liveness check |
| `POST /api/v1/documents` | Queue inline-content ingestion → **200** + job ID |
| `PUT /api/v1/documents` | Chunked async **file** ingestion (raw body) -> **202** + job ID |
| `GET /api/v1/documents` | List ingested documents |
| `DELETE /api/v1/documents/{id}` | Delete a document and all its chunks |
| `POST /api/v1/query` | Ask a grounded question |
| `GET /api/v1/jobs/{id}` | **Poll** ingestion job status/progress |

### `GET /health`

```bash
curl http://localhost:8080/health
# { "status": "ok" }
```

### `POST /api/v1/documents` — queue inline ingestion (returns immediately)

Validates the body, queues background ingestion, and responds **200** right away.

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Go is an open-source programming language designed at Google. It has built-in concurrency via goroutines and channels.",
    "metadata": { "source": "go-intro", "title": "About Go" }
  }'
```
```json
{ "job_id": "a1b2c3...", "status": "queued", "poll_url": "/api/v1/jobs/a1b2c3..." }
```

### `PUT /api/v1/documents` - chunked async file loading

Streams a file as the raw request body and ingests it in batches in the background.
Query parameters become document metadata. Responds **202 Accepted** with a job ID.

```bash
curl -X PUT "http://localhost:8080/api/v1/documents?source=handbook.txt&title=Handbook" \
  -H "Content-Type: text/plain" \
  --data-binary @handbook.txt
```
```json
{ "job_id": "d4e5f6...", "status": "queued", "poll_url": "/api/v1/jobs/d4e5f6..." }
```

### `GET /api/v1/jobs/{id}` - poll ingestion progress

```bash
curl http://localhost:8080/api/v1/jobs/d4e5f6...
```
```json
{
  "id": "d4e5f6...",
  "status": "running",
  "total_chunks": 128,
  "processed_chunks": 64,
  "document_id": "9f8e7d...",
  "created_at": "2026-06-29T12:00:00Z",
  "updated_at": "2026-06-29T12:00:03Z"
}
```

`status` transitions `queued -> running -> completed` (or `failed`, with an `error` field).
When complete, `document_id` is the ID to use for listing/deletion and appears as a
query source.

### `GET /api/v1/documents` - list documents

```bash
curl http://localhost:8080/api/v1/documents
```
```json
{
  "count": 1,
  "documents": [
    { "id": "9f8e7d...", "chunk_count": 128, "metadata": { "source": "handbook.txt", "title": "Handbook" } }
  ]
}
```

### `DELETE /api/v1/documents/{id}` - delete a document

```bash
curl -X DELETE http://localhost:8080/api/v1/documents/9f8e7d...
```
```json
{ "document_id": "9f8e7d...", "deleted_chunks": 128 }
```

### `POST /api/v1/query` - ask a grounded question

```bash
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{ "question": "How does Go handle concurrency?", "top_k": 3 }'
```
```json
{
  "answer": "Go provides built-in concurrency through goroutines and channels.",
  "sources": [
    {
      "id": "chunk-id...",
      "document_id": "9f8e7d...",
      "content": "Go is an open-source programming language ... goroutines and channels.",
      "score": 0.87,
      "metadata": { "source": "go-intro", "title": "About Go" }
    }
  ]
}
```

`top_k` is optional and defaults to `RETRIEVAL_TOP_K`.

### Error responses

Errors use a consistent envelope and meaningful status codes:

```json
{ "error": "document not found" }
```

| Status | When |
|--------|------|
| `400 Bad Request` | invalid JSON, empty content/question |
| `404 Not Found` | unknown job, unknown document, or no relevant documents for a query |
| `413 Payload Too Large` | uploaded file exceeds the limit (32 MiB) |
| `500 Internal Server Error` | upstream/model/store failure |

---

## How RAG works here

**Ingestion** (`POST` inline or `PUT` file) - asynchronous, with a job tracking progress:

```
Submit ──► create Job (queued) ──► return job id immediately
   background worker:
      TextSplitter.Split ──► for each batch: Embedder.EmbedDocuments ──► VectorStore.Add
                          └─► update Job progress (processed/total) ──► Job completed
```

**Query** (`POST /api/v1/query`):

```
question ──► Embedder.EmbedQuery ──► VectorStore.Search(topK)
          ──► PromptBuilder(question, sources) ──► LLM.Generate ──► answer + sources
```

The prompt instructs the model to answer **only** from the retrieved context and to
say it doesn't know otherwise — reducing hallucination and keeping answers grounded.

> **Note on persistence:** the default vector store and job store are in-memory, so
> ingested data and job history reset on restart. Implement `port.VectorStore` /
> `port.JobRepository` with a database to persist (see *Extending*).

---

## Extending

Because every dependency is a port, swapping an implementation is local and safe:

- **Persistent vector store** - implement `port.VectorStore` with
  [`pgvector`](https://github.com/pgvector/pgvector) and inject it in `cmd/server/main.go`.
- **Durable jobs** — implement `port.JobRepository` with Redis/Postgres.
- **Different embedder/LLM provider** - implement `port.Embedder` / `port.LLM`
  (e.g. Ollama, Cohere) and wire it in the composition root.
- **Custom prompting** - pass a different `usecase.PromptBuilder` to `NewQueryUseCase`.
- **Document loaders** - add a use case that uses langchaingo's `documentloaders`
  (PDF, HTML, CSV) ahead of the existing split -> embed -> store pipeline.

---

## Testing

```bash
make test            # or: go test ./...
```

The use cases are tested with in-memory fakes for every port, so the suite runs
**without any network access or API key** - a direct payoff of the dependency inversion.

---

## Docker

```bash
# Build and run with compose (reads variables from .env)
make docker-up

# Or build/run the image directly
make docker-build
docker run --rm -p 8080:8080 --env-file .env rag-server:latest
```

---
