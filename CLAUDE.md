# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run services
go run ./cmd/api-loyalty/main.go
go run ./cmd/points-engine/main.go

# Build binaries
go build -o api-loyalty ./cmd/api-loyalty
go build -o points-engine ./cmd/points-engine

# Start PostgreSQL (required for api-loyalty)
docker compose up -d

# Run tests
go test ./...
```

## Architecture

Two Go microservices under `cmd/`:

- **api-loyalty** (`cmd/api-loyalty/`) — REST API on port 8080. Connects to PostgreSQL via `pgx/v5`. Currently exposes `POST /transactions`. Database DSN is hardcoded.
- **points-engine** (`cmd/points-engine/`) — Standalone points calculation service. Currently a stub.

Module: `github.com/Quicksand06/loyalty` (Go 1.25, `pgx/v5` for Postgres).

Docker Compose provides PostgreSQL 16 (`loyalty` database, user `postgres`, password `Pa55word!`, port 5432). The app service definition is commented out — run services locally with `go run`.

No schema migrations or test files exist yet.
