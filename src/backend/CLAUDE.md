# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run backend locally (from repo root, loads .env)
task run-backend

# Run via Docker with hot-reload (Air)
task up

# Build binary
task build-backend   # outputs ./tmp/backend

# Run tests
cd src/backend && go test ./...

# Run a single test
cd src/backend && go test ./repository/ -run TestCleanupExpired_DeletesExpiredEntries
```

## Architecture

Layered architecture with manual dependency injection. All layers depend on interfaces, not concrete types.

```
main.go → di/ → handler/ → service/model/ → repository/
                              ↑                    ↑
                        service/interface/     domain/
```

- **main.go** — Wires everything: loads config, creates S3 client, initializes SQLite, calls DI factories, starts Gin server on `:8080`, launches background cache cleanup goroutine
- **di/** — One factory function per handler (`CreateXHandler()`). Manual wiring, no framework
- **handler/** — Gin HTTP handlers. Each handler depends on a service interface. All endpoints under `/api/v1`
- **service/interface/** — Defines all interfaces: services, repositories. This is the contract layer
- **service/model/** — Business logic implementations. Cache coordination logic lives here (not in repositories)
- **repository/** — Data access: S3 via AWS SDK, SQLite for cache entries and settings, go-cache for in-memory list cache
- **domain/** — Pure DTOs with JSON tags. No behavior
- **config/** — Loads R2 credentials and cache settings from environment variables
- **infrastructure/** — SQLite setup with WAL mode and schema migration

## Key Patterns

**Handler pattern**: Extract params from `*gin.Context`, call service, return JSON (200) or error (500).

**Multi-tier caching**:

- In-memory (go-cache): bucket lists (60min TTL), object lists (10min TTL, per bucket+prefix)
- Disk + SQLite: file content cache with configurable TTL and max size. ETag-aware invalidation. Background cleanup goroutine with periodic eviction

**CacheRepository singleton**: Uses `sync.Once`. Tests reset the singleton via `once = sync.Once{}` before each test.

**Atomic file writes**: Content cache writes to `.tmp` then renames to final path.

## Environment Variables

| Variable                         | Default           | Description                                    |
| -------------------------------- | ----------------- | ---------------------------------------------- |
| `R2_ACCOUNT_ID`                  | (required)        | Cloudflare R2 account ID                       |
| `R2_ACCESS_KEY_ID`               | (required)        | R2 access key                                  |
| `R2_SECRET_ACCESS_KEY`           | (required)        | R2 secret key                                  |
| `env`                            | -                 | Set `dev` to trust local network IPs           |
| `IP_LIST`                        | -                 | Comma-separated trusted proxy IPs (production) |
| `CACHE_DB_PATH`                  | `./data/cache.db` | SQLite database path                           |
| `CACHE_DIR`                      | `./data/cache`    | Cached file storage directory                  |
| `CACHE_TTL_MINUTES`              | `120`             | File cache TTL                                 |
| `CACHE_CLEANUP_INTERVAL_MINUTES` | `60`              | Background cleanup interval                    |
| `CACHE_MAX_SIZE_MB`              | unlimited         | Max disk cache size (LRU eviction)             |

## 補足

ユーザーへの対話は日本語を利用してください
