# syntax=docker/dockerfile:1

# --- build stage ---
FROM golang:1.23-alpine AS build
WORKDIR /src

# Cache dependencies first (go.sum is optional and copied when present).
COPY go.mod go.sum* ./
RUN go mod download

# Build a static binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/rag-server ./cmd/server

# --- runtime stage ---
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/rag-server /app/rag-server
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/rag-server"]
