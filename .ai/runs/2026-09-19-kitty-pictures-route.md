# Add `/kitty` random colour-aware pictures

## Goal
Add a `GET /kitty` page that displays exactly one randomly selected embedded kitty picture per request, using a black kitty during even-numbered minutes and a red kitty during odd-numbered minutes.

## Scope
- Add a small, testable kitty selection helper that derives the required colour from the request-time minute and randomly selects one image from that colour's embedded assets.
- Add the `/kitty` route, handler, page template, and embedded SVG kitty assets.
- Add handler and selection tests covering route behavior, one-image rendering, minute parity, and invalid paths.

## Non-goals
- No external image service or runtime network request.
- No changes to the existing home page, middleware, server configuration, or unrelated asset behavior.
- No persistent storage, user preferences, or new dependencies.

## Risks
- A template or asset path mismatch would turn a valid request into a server error; tests will exercise rendered output and build-time embedding.
- The configured Go validation gate is currently blocked by sandbox permissions reading the Go toolchain cache under the user home directory; rerun it after the environment allows that access.
- Randomness must select exactly one asset without allowing a colour mismatch; the helper will constrain the candidate set before selection.

## Implementation Plan

### Phase 1: Kitty selection and assets

1.1 Add embedded black and red kitty SVG assets and a selection helper with injectable time/randomness for deterministic tests.

### Phase 2: HTTP surface and verification

2.1 Add the `/kitty` handler, template, and route while preserving existing middleware and route behavior.
2.2 Add focused tests for parity, random selection, one-image output, and route errors; run the configured validation gate.

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles.

### Phase 1: Kitty selection and assets

- [x] 1.1 Add embedded black and red kitty SVG assets and a selection helper with injectable time/randomness for deterministic tests. — 5c1f096

### Phase 2: HTTP surface and verification

- [x] 2.1 Add the `/kitty` handler, template, and route while preserving existing middleware and route behavior. — 2903f05
- [ ] 2.2 Add focused tests for parity, random selection, one-image output, and route errors; run the configured validation gate.
