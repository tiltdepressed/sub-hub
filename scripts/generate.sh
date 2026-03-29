#!/usr/bin/env bash
set -euo pipefail

echo "==> sqlc (if installed)"
if command -v sqlc >/dev/null 2>&1; then
	(cd internal/db/sql && sqlc generate)
else
	echo "sqlc not found; skip"
fi

echo "==> oapi-codegen (if installed)"
if command -v oapi-codegen >/dev/null 2>&1; then
	mkdir -p internal/generated/openapi
	oapi-codegen -generate types -package openapi api/openapi.yaml > internal/generated/openapi/types.gen.go
	oapi-codegen -generate chi-server -package openapi api/openapi.yaml > internal/generated/openapi/server.gen.go
else
	echo "oapi-codegen not found; skip"
fi

