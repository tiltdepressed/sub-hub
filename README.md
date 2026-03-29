# sub-hub

REST service for managing subscriptions and calculating total cost for a period.

Stack: Go (chi), PostgreSQL (pgx), goose migrations, zap logging, Prometheus metrics, OpenAPI + Swagger UI, optional OTEL tracing.

## Quick start (local)

Requirements: Go `1.25.x`.

```bash
cp .env.example .env
make tidy
make run
```

Service will be available at `http://localhost:8080`.

## Quick start (Docker Compose)

```bash
make docker-up
```

What you get:
- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger/`
- Prometheus: `http://localhost:9090`
- Jaeger UI: `http://localhost:16686`

Postgres is mapped to the host as `localhost:5433` (inside the compose network it is `postgres:5432`).

Stop:

```bash
make docker-down
```

## Migrations

Requires `goose` installed and `DB_DSN` set.

```bash
export DB_DSN='postgres://postgres:postgres@localhost:5432/subhub?sslmode=disable'
./scripts/migrate.sh
```

If Postgres is running via compose and you run the app locally, use port `5433`.

## Configuration

Main environment variables (see `.env.example`):
- `HTTP_ADDR` — HTTP server address, e.g. `:8080`
- `DB_DSN` — Postgres connection string
- `LOG_LEVEL` — log level
- `TRACING_ENABLED` — enable OTEL (`true/false`)
- `TRACING_ENDPOINT` — OTLP gRPC endpoint for tracing
- `TRACING_SERVICE_NAME` — tracing service name
- `OPENAPI_PATH` — path to `api/openapi.yaml`

## API

Service endpoints:
- `GET /healthz`
- `GET /readyz`
- `GET /metrics`
- `GET /swagger/`
- `GET /api/v1/ping`

Subscriptions:
- `POST /api/v1/subscriptions`
- `GET /api/v1/subscriptions`
- `GET /api/v1/subscriptions/{id}`
- `PUT /api/v1/subscriptions/{id}`
- `DELETE /api/v1/subscriptions/{id}`
- `GET /api/v1/subscriptions/total?from=MM-YYYY&to=MM-YYYY[&user_id=...][&service_name=...]`

Date format: `MM-YYYY` (e.g. `07-2025`).

## Example requests

Create a subscription:

```bash
curl -fsS -X POST 'http://localhost:8080/api/v1/subscriptions' \
	-H 'Content-Type: application/json' \
	-d '{
		"service_name": "Yandex Plus",
		"price": 400,
		"user_id": "00000000-0000-0000-0000-000000000001",
		"start_date": "07-2025",
		"end_date": "10-2025"
	}'
```

Get total cost for a period:

```bash
curl -fsS 'http://localhost:8080/api/v1/subscriptions/total?from=07-2025&to=09-2025&user_id=00000000-0000-0000-0000-000000000001'
```

## Dev commands

```bash
make test
make build
make run
```

Linter (if `golangci-lint` is installed):

```bash
make lint
```

Generation (if `sqlc` and/or `oapi-codegen` are installed):

```bash
make generate
```
