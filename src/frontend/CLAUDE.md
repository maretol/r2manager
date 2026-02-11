# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
npm run dev      # Start Next.js dev server on :3000
npm run build    # Build standalone output
npm run lint     # ESLint (core-web-vitals + typescript configs)
```

No tests exist in this project.

## Architecture

Next.js 16 App Router frontend for the R2 Manager application. Communicates with a Go backend API via a catch-all API proxy route.

### API Proxy Pattern

`app/api/v1/[...path]/route.ts` proxies all `/api/v1/*` requests to the Go backend (`BACKEND_URL`). All backend communication from pages uses `lib/api.ts`, which builds URLs with `SERVER_URL` — this resolves to the local Next.js server (self-referencing via `HOSTNAME`) on the server side, and to a relative path on the client side.

### Server vs Client Components

- **Pages** (`app/**/page.tsx`) are async Server Components that fetch data via `lib/api.ts` and pass it as props to child components
- **Interactive components** use `'use client'` directive — `object-detail-panel.tsx`, `settings-form.tsx`, `refresh-*-button.tsx`, `bucket-menu-item.tsx`
- **Server Actions** in `app/**/actions.ts` handle mutations (object URL generation, settings save)

### Key Patterns

- Object selection uses URL search params (`?selected=key&prefix=path/`) rather than client state — the page server component resolves the selected object
- `lib/object-utils.ts` transforms flat R2 object lists into a folder/file display hierarchy (`DisplayObject[]`) with folders sorted first
- Date formatting in `object-utils.ts` manually converts UTC to JST (+9h) rather than using locale APIs
- `components/ui/` contains shadcn/ui (new-york style, Radix-based) — do not manually edit these files; use `npx shadcn@latest add <component>`

### Environment Variables

- `BASE_PATH` — URL prefix for reverse proxy deployment (used in `next.config.ts` and server actions)
- `BACKEND_URL` — Go backend URL (used only by the API proxy route, default `http://localhost:8080`)
- `HOSTNAME`, `PORT`, `PROTOCOL` — Used by `lib/api.ts` to construct server-side self-referencing URLs

## Code Style

- Prettier: single quotes, no semicolons, 2-space indent, 120 char width
- Use `@/*` path aliases for all imports (maps to project root)
- Style with Tailwind CSS v4
- Icons from `lucide-react`
- UI language: English for UI labels, Japanese for user-facing error/status messages and code comments

## 補足

ユーザーへの対話は日本語を利用してください
