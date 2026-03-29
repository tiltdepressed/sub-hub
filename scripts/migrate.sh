#!/usr/bin/env bash
set -euo pipefail

: "${DB_DSN:?DB_DSN is required}"

if ! command -v goose >/dev/null 2>&1; then
	echo "goose not found"
	exit 1
fi

goose -dir internal/db/migrations postgres "$DB_DSN" up

