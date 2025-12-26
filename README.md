# shopping

Small web app to track supplies: what’s missing, current levels, and grouping (e.g. carrots → vegetables).

## Requirements

- Go 1.25.5+
- SQLite
- `air` (optional, for live reload)
- `migrate` (golang-migrate CLI, optional but recommended)

## Configure

Copy and edit:

```bash
cp configs/dev.env.example .env
set -a; source .env; set +a
```

Key env vars:

- `ADDR` (default `:8080`)
- `DB_DSN` (default `file:data/shopping.db?cache=shared&mode=rwc&_pragma=foreign_keys(1)`)
- `BASE_URL` (e.g. `http://localhost:8080`)
- OIDC: `OIDC_ISSUER`, `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `OIDC_REDIRECT_URL`
- Dev-only bypass: `AUTH_DISABLED=1`

## Database migrations (golang-migrate)

Migrations live in `migrations/`.

On startup, the app automatically applies embedded migrations.

Example:

```bash
mkdir -p data
migrate -path ./migrations -database "sqlite3://data/shopping.db" up
```

## Run

```bash
make run
```

With `air`:

```bash
make dev
```

## Make targets

```bash
make run
make dev
make migrate-up
```

## Admin

- `POST /admin/db/optimize` runs `PRAGMA optimize`
- If OIDC token includes claim `admin=true`, admin access is granted.
- Otherwise requires `ADMIN_EMAILS` (comma-separated) unless `AUTH_DISABLED=1`

## sqlc

SQL queries live in `internal/db/queries.sql` and are intended to be compiled with `sqlc` using `sqlc.yaml`.

## CQRS

- Reads (`InventoryQueries`) query SQL views (`v_groups`, `v_products`) created by migrations.
- Writes (`InventoryCommands`) operate on base tables (`groups`, `products`).

## Architecture

- Domain (pure logic + fast unit tests): `internal/domain/products`, `internal/domain/admin`
- Infrastructure:
  - Config: `internal/infrastructure/config`
  - OIDC auth: `internal/infrastructure/oidc`
  - DB/persistence: `internal/infrastructure/persistence/sqlite`
  - Logging: `internal/infrastructure/logging`
- Web UI: `internal/web` renders HTML with `github.com/a-h/templ` and uses htmx for interactions.
