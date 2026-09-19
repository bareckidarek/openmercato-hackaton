# Code review rules

## Review priorities

1. **Correctness:** verify HTTP behavior, status codes, method/path handling, graceful shutdown, and error paths.
2. **Security:** preserve security headers, panic recovery, request boundaries, HTML escaping, and safe handling of environment/configuration input.
3. **Contracts:** treat routes, response formats, embedded template names, command-line flags, and environment variable defaults as compatibility surfaces.
4. **Maintainability:** prefer the existing `application` methods, `internal/request`, `internal/response`, and `internal/validator` helpers over parallel implementations.

## Repository-specific checks

- Run `make test` and `make build`; use `make audit` for broad or security-sensitive changes.
- Keep Go code formatted with `gofmt`; inspect embedded assets when changing template or static-file paths.
- For handler changes, check the happy path, wrong path/method behavior, server errors, and relevant middleware ordering.
- For configuration changes, verify environment parsing, defaults, invalid values, and the `--version` flag.
- For template changes, verify escaped output and the default/empty/error states that the page can display.
- Review dependency or module changes for unnecessary additions and run `go mod verify`.

## Severity guidance

- **Blocker:** security vulnerability, data loss, broken build, or a breaking contract without an approved migration path.
- **Major:** incorrect user-visible behavior, missing validation/error handling, or a regression in routing or shutdown.
- **Minor:** localized maintainability, test coverage, or documentation issue that should be addressed before merge when practical.
- **Nit:** non-blocking style preference.
