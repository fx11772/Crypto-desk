#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")/.."

COMPOSE=(podman compose -f compose.test.yaml -p flux-test)

cleanup() {
    echo "Cleaning up integration test environment..."
    "${COMPOSE[@]}" down --volumes
}

trap cleanup EXIT

echo "Removing previous test environment..."
cleanup

echo "Starting PostgreSQL test database..."
"${COMPOSE[@]}" up -d flux-db-test

echo "Applying database migrations..."
"${COMPOSE[@]}" run --rm flux-migrate-test up

echo "Running integration tests..."
"${COMPOSE[@]}" run --rm flux-tests