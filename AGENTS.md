# Repository Guidelines

## Project Structure & Module Organization

- `cmd/shopping/`: main entrypoint (`go run ./cmd/shopping`).
- `internal/domain/`: pure domain code intended for fast unit tests.
  - `internal/domain/products/`: product/group entities, CQRS ports, validation helpers.
  - `internal/domain/admin/`: admin ports (e.g., DB optimize).
- `internal/domain/products/`: domain entities + validation + write-side service (`Service`) and repo port (`Repository`).
- `internal/infrastructure/`: adapters and integrations.
  - `internal/infrastructure/persistence/sqlite/`: SQLite + `sqlc` repo (implements domain ports).
  - `internal/infrastructure/oidc/`: OIDC auth against Tailscale `tsidp`.
  - `internal/infrastructure/config/`: env + `.env` loading.
  - `internal/infrastructure/logging/`: `slog` logger setup.
- `internal/web/`: HTTP handlers and templates wiring (controllers should stay thin).
- `web/`: static assets (htmx-driven UI).
- Views are rendered in Go using `github.com/a-h/templ` (see `internal/web/views_templ.go`).
- `migrations/`: SQL migrations (also embedded and run on startup).
- `internal/db/`: `sqlc`-generated query code; edit `internal/db/queries.sql`, not `*.gen.go`.

## Build, Test, and Development

- `make build`: builds `bin/shopping`.
- `make run`: runs locally (loads `.env` automatically).
- `make dev`: runs with live reload via `air`.
- `make test`: runs `go test ./...`.
- `make fmt`: runs `gofmt` on Go sources.
- `make migrate-up`: applies migrations using `migrate` CLI (optional; app also auto-migrates).

## Coding Style & Naming Conventions

- Go formatting: `gofmt` is required (tabs, standard Go style).
- Packages: prefer domain-first names (`products`, `admin`) over generic “service” packages.
- Generated code: keep `internal/db/*.gen.go` and `internal/db/models.go` unchanged (regenerate via `sqlc`).

## Testing Guidelines

- Use Go’s standard `testing` package.
- Keep domain tests fast and isolated (no DB/network). Example: `internal/domain/products/validation_test.go`.
- Run all tests with `make test`.

## Commit & Pull Request Guidelines

- Commit messages: use a short imperative summary (e.g., `Add OIDC callback handling`, `Refactor products domain ports`).
- PRs should include: intent/summary, how to test (`make test`, `make run`), and any config changes (new env vars, migrations).

## Security & Configuration Tips

- Never commit secrets; `.env` is ignored. Rotate OIDC secrets if exposed.
- Admin endpoints require OIDC claim `admin=true` or `ADMIN_EMAILS` (unless `AUTH_DISABLED=1`).
