# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

R2 Manager is a web application for managing files stored in Cloudflare R2 (S3-compatible storage). It is designed for personal use on a self-hosted Docker environment. Authentication/authorization is delegated to a reverse proxy (nginx); the app itself has no auth. Multi-user support is not planned.

Review comments and commit messages should be written in Japanese.

## Commands

### Development (Docker, recommended)

```bash
task up        # Start dev containers (backend + frontend + nginx)
task down      # Stop dev containers
```

Dev environment: nginx on `:80` → frontend (Next.js dev) on `:3000`, backend (Go with Air hot-reload) on `:8080`.

### Development (local, without Docker)

```bash
task run-backend    # Run Go backend (reads .env from root)
task run-frontend   # Run Next.js frontend (reads .env from root)
```

### Build

```bash
task build-backend                          # Build Go binary to ./tmp/backend
cd src/frontend && npm run build            # Build Next.js standalone output
task up-prd                                 # Build and start production containers
```

### Linting

```bash
cd src/frontend && npm run lint   # ESLint for frontend
```

### No tests

There are currently no tests in this project (neither Go tests nor frontend tests).

## Architecture

**Two-service architecture**: Go backend API + Next.js frontend, communicating via REST.

### Backend (`src/backend/`) — Go + Gin

Layered architecture with dependency injection:

- `main.go` — Entry point: loads config, initializes S3 client, SQLite DB, DI container, starts Gin server on `:8080`
- `handler/` — HTTP handlers (Gin)
- `router/` — Route definitions. All API routes under `/api/v1`
- `service/` — Business logic. `service/interface/` defines interfaces, implementations are in `service/`
- `repository/` — Data access (R2 via AWS SDK, SQLite for cache/settings)
- `domain/` — Domain models (Bucket, Object, BucketSettings, CacheEntry)
- `di/` — Dependency injection wiring
- `config/` — Configuration loading (R2 credentials, cache settings)
- `infrastructure/` — Database setup (SQLite with WAL mode)

Key API endpoints:

- `GET /api/v1/buckets` — List buckets
- `GET /api/v1/buckets/:bucketName/objects` — List objects (paginated)
- `GET /api/v1/buckets/:bucketName/content/*key` — Stream content from R2
- `GET/PUT /api/v1/settings/buckets` — Bucket settings (public URL config)
- `DELETE /api/v1/cache/content`, `/api/v1/cache/api` — Cache management

### Frontend (`src/frontend/`) — Next.js 16 App Router + TypeScript

- `app/` — Pages using App Router. Server Components fetch data; Client Components handle interactivity
- `app/api/v1/[...path]/route.ts` — API proxy that forwards all `/api/v1/*` requests to the Go backend
- `components/` — React components. `components/ui/` contains shadcn/ui (Radix-based) components
- `lib/` — Utility functions and API client helpers
- `hooks/` — Custom React hooks
- `types/` — TypeScript type definitions

### Environment Variables

**Backend** (`src/backend/.env`): `R2_ACCOUNT_ID`, `R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY`, `env` (dev/prd), `IP_LIST`

**Frontend** (`src/frontend/.env.local`): `BASE_PATH` (URL prefix for reverse proxy), `BACKEND_URL` (backend API URL), `NEXT_PUBLIC_APP_NAME`

## Code Style

- **Frontend**: Prettier with single quotes, no semicolons, 2-space indent, 120 char width. Use `@/*` path aliases for imports. Style with Tailwind CSS.
- **Backend**: Standard Go conventions.

## 補足

ユーザー側には日本語でコメントしてください
