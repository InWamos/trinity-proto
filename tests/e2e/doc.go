package e2e

// Package e2e contains end-to-end tests for the Trinity API.
//
// These tests use testcontainers to spin up real instances of:
// - PostgreSQL (postgres:18.1-trixie)
// - Redis (redis:8.4-bookworm)
// - Migrate (migrate/migrate:v4.19.1)
//
// Run tests with:
//   go test -v ./tests/e2e/...
//
// Note: Docker must be running and accessible.
