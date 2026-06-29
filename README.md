# RAG Backend (Go · Clean Architecture · langchaingo)

A production-style **Retrieval-Augmented Generation (RAG)** backend written in Go.
It exposes a small HTTP API to **ingest documents** and **ask grounded questions**:
documents are chunked, embedded, and stored as vectors; questions retrieve the most
relevant chunks and feed them to an LLM that answers **only** from that context.

Built with [**langchaingo**](https://github.com/tmc/langchaingo) using the **OpenAI**
provider for both chat completion and embeddings, and structured around
**Clean Architecture** with **SOLID**, **DRY**, and **KISS** as first-class constraints.

---

## Features

- 🧩 **Clean Architecture** — domain, application, infrastructure, and delivery layers with a strict inward dependency rule.
- 🔌 **Pluggable by design** — LLM, embedder, vector store, and text splitter are interfaces (ports); swap implementations without touching business logic.
- 🤖 **OpenAI via langchaingo** — any OpenAI chat model (configurable through `LLM_MODEL`) plus `text-embedding-3-small` embeddings.
- 🗃️ **Zero-setup vector store** — dependency-free in-memory cosine-similarity store; perfect for local dev and tests, trivially replaceable with pgvector/Qdrant.
- 🌐 **Standard-library HTTP** — method-aware routing (Go 1.22+ `ServeMux`), structured logging (`log/slog`), panic recovery, and graceful shutdown.
- 🧪 **Tested core** — the use cases and store are unit-tested with fakes (no network required).
- 🐳 **Container-ready** — multi-stage `Dockerfile` (distroless) and `docker-compose.yml`.

---

## Architecture

The codebase follows **Clean Architecture**: dependencies point **inward only**.
Inner layers know nothing about outer layers — the domain has zero imports from
infrastructure or transport.

```
┌──────────────────────────────────────────────────────────────┐
│  Delivery (HTTP)            cmd/server, internal/delivery/rest │  ← frameworks & I/O
│   router · handlers · middleware · DTOs                        │
│        │ depends on                                            │
│        ▼                                                       │
│  Application (Use Cases)    internal/usecase                   │  ← business workflows
│   IngestUseCase · QueryUseCase · PromptBuilder                 │
│        │ depends on                                            │
│        ▼                                                       │
│  Domain (Core)              internal/domain                    │  ← entities & ports
│   entities · ports (LLM, Embedder, VectorStore, TextSplitter) │
│        ▲ implemented by                                        │
│        │                                                       │
│  Infrastructure (Adapters)  internal/infrastructure           │  ← langchaingo, stores
│   OpenAI LLM · OpenAI Embedder · Splitter · In-memory store   │
└──────────────────────────────────────────────────────────────┘
```

**The Dependency Rule.** `domain` imports nothing from the project. `usecase` imports
only `domain`. `infrastructure` implements `domain` ports. `delivery` depends on
`usecase` interfaces. The **composition root** (`cmd/server/main.go`) is the only place
that knows every concrete type — it wires them together.

### How the principles map to the code

| Principle | Where you see it |
|-----------|------------------|
| **S**ingle Responsibility | Each type does one thing: `ingestInteractor` orchestrates ingestion, `Memory` stores vectors, `RAGHandler` translates HTTP ↔ use cases. |
| **O**pen/Closed | Add a pgvector store or a different LLM by implementing a port — no existing code changes. The `var _ port.X = (*T)(nil)` compile-time checks guard the contracts. |
| **L**iskov Substitution | Any `port.VectorStore` (in-memory today, pgvector tomorrow) is interchangeable; the use cases are oblivious. |
| **I**nterface Segregation | Small, focused ports (`LLM`, `Embedder`, `VectorStore`, `TextSplitter`) instead of one fat "AI service" interface. |
| **D**ependency Inversion | Use cases and handlers depend on interfaces; concretes are injected at the composition root. |
| **DRY** | Shared JSON/error helpers, one prompt builder, one config loader, reusable middleware chain. |
| **KISS** | Standard-library HTTP and routing, an in-memory store, env-based config — no framework ceremony. |

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
│   │   ├── entity/                  # Document, ScoredDocument, Answer
│   │   ├── port/                    # LLM, Embedder, VectorStore, TextSplitter
│   │   └── errors.go                # sentinel domain errors
│   ├── usecase/                     # application business rules
│   │   ├── ingest.go                # split → embed → store
│   │   ├── query.go                 # embed → retrieve → prompt → generate
│   │   └── prompt.go                # grounded RAG prompt builder
│   ├── infrastructure/              # adapters implementing the ports
│   │   ├── llm/openai.go            # langchaingo OpenAI chat
│   │   ├── embedding/openai.go      # langchaingo OpenAI embeddings
│   │   ├── splitter/recursive.go    # langchaingo recursive splitter
│   │   └── vectorstore/memory.go    # in-memory cosine store
│   └── delivery/
│       └── rest/                    # HTTP transport
│           ├── router.go            # routes + middleware
│           ├── handler/             # ingest, query, health
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
- **langchaingo** — LLM, embeddings, and text-splitting abstractions
- **OpenAI** — chat completion (`gpt-4o-mini` by default) and `text-embedding-3-small`
- **godotenv** — `.env` loading for local development

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
| `OPENAI_API_KEY` | — | **Required.** OpenAI API key |
| `OPENAI_BASE_URL` | _(unset)_ | Override base URL (Azure OpenAI / proxy / compatible server) |
| `LLM_MODEL` | `gpt-4o-mini` | Any OpenAI chat model (`gpt-4o`, `gpt-4-turbo`, `gpt-3.5-turbo`, …) |
| `LLM_TEMPERATURE` | `0.2` | Sampling temperature for generation |
| `EMBEDDING_MODEL` | `text-embedding-3-small` | OpenAI embedding model |
| `CHUNK_SIZE` | `1000` | Characters per chunk |
| `CHUNK_OVERLAP` | `200` | Overlap between adjacent chunks |
| `RETRIEVAL_TOP_K` | `4` | Number of chunks retrieved per query |

---

## API reference

### `GET /health`

```bash
curl http://localhost:8080/health
```
```json
{ "status": "ok" }
```

### `POST /api/v1/documents` — ingest a document

Splits the content into chunks, embeds them, and stores the vectors.

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Go is an open-source programming language designed at Google. It is statically typed and compiled, with built-in concurrency via goroutines and channels.",
    "metadata": { "source": "go-intro", "title": "About Go" }
  }'
```
```json
{ "chunks_created": 1 }
```

### `POST /api/v1/query` — ask a grounded question

Embeds the question, retrieves the most similar chunks, and asks the LLM to answer
using only that context. Returns the answer plus the source chunks (with scores).

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
      "id": "a1b2c3d4...",
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
{ "error": "question must not be empty" }
```

| Status | When |
|--------|------|
| `400 Bad Request` | invalid JSON, empty content/question |
| `404 Not Found` | no relevant documents for the query |
| `500 Internal Server Error` | upstream/model/store failure |

---

## How RAG works here

**Ingestion** (`POST /api/v1/documents`):

```
content ──► TextSplitter.Split ──► Embedder.EmbedDocuments ──► VectorStore.Add
```

**Query** (`POST /api/v1/query`):

```
question ──► Embedder.EmbedQuery ──► VectorStore.Search(topK)
          ──► PromptBuilder(question, sources) ──► LLM.Generate ──► answer + sources
```

The prompt instructs the model to answer **only** from the retrieved context and to
say it doesn't know otherwise — reducing hallucination and keeping answers grounded.

---

## Extending

Because every dependency is a port, swapping an implementation is local and safe:

- **Persistent vector store** — implement `port.VectorStore` with
  [`pgvector`](https://github.com/pgvector/pgvector) (langchaingo has a `vectorstores/pgvector`
  package) and inject it in `cmd/server/main.go`. Nothing else changes.
- **Different embedder/LLM provider** — implement `port.Embedder` / `port.LLM`
  (e.g. Ollama, Cohere) and wire it in the composition root.
- **Custom prompting** — pass a different `usecase.PromptBuilder` to `NewQueryUseCase`.
- **Document loaders** — add a use case that uses langchaingo's `documentloaders`
  (PDF, HTML, CSV) before the existing split → embed → store pipeline.

---

## Testing

```bash
make test            # or: go test ./...
```

The use cases are tested with in-memory fakes for every port, so the suite runs
**without any network access or API key** — a direct payoff of the dependency inversion.

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

## License

[MIT](LICENSE) © Nurzhek
