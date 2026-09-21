# AGENTS.md

WindNav (微风导航): a self-hosted navigation/homepage app. Frontend/backend split in one repo.

## Layout

- `backend/` — Go 1.26, Gin + GORM. Entrypoint `cmd/server/main.go`. Deps are **vendored** (`backend/vendor/`), so build with `-mod=vendor`.
- `frontend/` — Vue 3 + Vite + TypeScript + Tailwind CSS v4 + Pinia + TanStack Query. Entrypoint `src/main.ts`.
- `deploy/` — Dockerfile + two compose files (SQLite single-container, PostgreSQL standard).
- `docs/` — Chinese-language run/debug and API docs. `plans/` — design docs (may be aspirational, not current code).

## Commands

Backend (run from `backend/`):

```bash
go run ./cmd/server
go vet ./...
go test ./...            # unit
go test -tags=integration ./...   # integration suite (CI runs this)
go build -mod=vendor ./cmd/server
```

Frontend (run from `frontend/`):

```bash
npm install
npm run dev      # Vite dev server on :5173
npm run build    # runs vue-tsc -b (typecheck) THEN vite build — typecheck failures fail the build
```

CI (`.gitea/workflows/windnav-dev.yaml`) runs: `go vet` → `go test` → integration tests → `go build -mod=vendor` for backend; `npm ci` + `npm run build` for frontend. Match this order locally.

## Config & environment

Backend config is **environment-variable only** — there is no `.env` file loading. Set vars in the shell before `go run` (see `.env.example`). Key vars: `DB_DRIVER` (`sqlite`|`postgres`), `DB_SQLITE_PATH`, `DATABASE_URL`, `JWT_SECRET`, `ADMIN_USERNAME`, `ADMIN_PASSWORD`, `HTTP_ADDR` (default `:8080`), `STATIC_DIR`.

- Default DB is SQLite at `./data/windnav.db`; the file, schema, and default admin are auto-created on startup (GORM AutoMigrate + Seed). `backend/data/` is gitignored.
- **Admin password is only seeded if the user does not already exist.** Changing `ADMIN_PASSWORD` after first boot does NOT update an existing admin — delete the DB file and restart to reset.
- Production (`APP_ENV=production`) refuses to start with the default `JWT_SECRET`.
- `STATIC_DIR` (e.g. `frontend/dist`) makes the backend serve the built frontend; otherwise the frontend is served separately.

## Local dev

- Backend: `:8080` (health check `GET /healthz`). Frontend: `:5173`.
- Vite proxies `/api` and `/healthz` to `localhost:8080` (see `frontend/vite.config.ts`). If you change backend `HTTP_ADDR`, update this proxy too.
- Default admin: `admin` / `admin123456`.
- Admin UI at `/admin`; public page at `/`.

## Deploy gotchas

- `deploy/docker-compose*.yml` map host port **`80:80`**, but the app listens on `:8080` (Dockerfile `ENV HTTP_ADDR=:8080`). The published port and app port don't match — don't "fix" one without the other.
- Dockerfile builds backend with `CGO_ENABLED=1` (required for the SQLite driver) and `-mod=vendor`.

## API conventions

- Unified response envelope: `{ "data": ..., "error": null, "meta": {} }`; errors use `{ "code", "message", "details" }`.
- Public routes: `/api/public/*` (no auth). Auth: `/api/auth/*`. Admin: `/api/admin/*` — all require `Authorization: Bearer <token>`.
- Login response: `data.token` (JWT) + `data.user`. Frontend stores it in `localStorage` (key in `src/api/client.ts`).

## Notes

- Repo docs are in Chinese; keep new docs consistent.
- `docs/` and `plans/` may lag the code — trust the code and CI as source of truth.