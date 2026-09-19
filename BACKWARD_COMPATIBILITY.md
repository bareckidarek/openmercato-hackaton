# Backward compatibility

This document records the public surfaces that must be considered before changing behavior. The application is not currently a published Go library and has no database migrations or event/webhook contracts in the repository.

## HTTP routes and responses

- `GET /` serves the home page. A non-root path reaching the home handler returns the existing not-found response.
- `GET /static/` serves embedded static assets.
- Middleware ordering, security headers, recovery behavior, and status codes are part of the HTTP contract.

A breaking change removes or changes an existing route, method, status code, response shape, or security behavior. Update callers and README guidance in the same change, and include tests covering the old and new behavior.

## CLI and configuration

- `go run ./cmd/web` starts the server.
- `--version` prints the application version and exits.
- `BASE_URL` and `HTTP_PORT` are supported environment variables with existing defaults.

A breaking change removes a flag or environment variable, changes its meaning, or changes a default that affects deployment. Document the migration and test the compatibility path before merging.

## Embedded assets and templates

`assets/templates/` and `assets/static/` are embedded into the binary. Renaming or removing a template or static path is a runtime contract change; update every reference and verify `make build` plus a request-level test.

## Go package APIs

Packages under `internal/` are repository-internal and are not importable by external modules under Go's `internal` rules. Changes still require updating in-repository callers and tests.
