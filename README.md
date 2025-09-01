# task_effective_mobile_subscribe

Small Go service to manage user subscriptions (HTTP API). This README explains local setup, configuration, swagger, migrations and developer tooling.

## Quick start

1. Install Go 1.25+ and ensure `$GOPATH/bin` is in your PATH.
2. Install required CLI tools (recommended pinned versions in `go.mod` / `tools.go`):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

3. Create a `.env` or export envs used by the app (example):

```bash
export PSQL_SOURCE="postgres://user:pass@localhost:5432/subs_db?sslmode=disable"
export PORT=8080
export LOG_LEVEL=DEBUG  # DEBUG|INFO|WARN|ERROR

```

4. Run database migrations (example using `goose`):

```bash
# from project root, adjust URL and path to migrations
goose  up
```

5. Generate swagger docs (if you change annotations):

```bash
swag init -g main.go -o docs
```

6. Build and run:

```bash
go build ./...
go run main.go
```

Open the API: http://localhost:${PORT}/api/subscriptions
Swagger UI: http://localhost:${PORT}/swagger/index.html

## Configuration

Configuration is loaded from environment variables. Important keys:

- `PSQL_SOURCE` — Postgres connection string
- `PORT` — HTTP listen port (e.g. `:8080` or `:3000`)
- `LOG_LEVEL` — logging level (debug/info/warn/error)

Add other external API URLs, keys and toggles to the config and avoid hardcoding them.


## Migrations

Migrations live in `internal/database/migrations`. Use `migrate` to apply/rollback. Example commands:

```bash
# apply all up
goose -dir internal/database/migrations postgres "$PSQL_SOURCE" up
# rollback one
goose -dir internal/database/migrations postgres "$PSQL_SOURCE" down 1
```
