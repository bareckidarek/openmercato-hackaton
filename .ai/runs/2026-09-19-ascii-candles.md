# Implement `/ascii-candles` spec

Source doc: .ai/specs/2026-09-19-ascii-candles.md

## Goal
Implement the reviewed ASCII SPY daily candlestick endpoint using validated Yahoo Finance data and the existing Go HTTP patterns.

## Scope
Add an injected market-data client, validated OHLC model, deterministic ASCII renderer, GET route, error handling, and unit/HTTP tests. Output must be strictly ASCII.

## Non-goals
No persistence, caching, configurable symbols/ranges, trading signals, frontend UI, or external charting library.

## Progress

> Convention: `- [ ]` pending, `- [x]` done. Append ` — <commit sha>` when a step lands. Do not rename step titles.

### Phase 1: Data and renderer

- [x] 1.1 Add validated Yahoo chart client and OHLC model with fake-server tests. — 707b1dc
- [x] 1.2 Add deterministic ASCII candle renderer with ASCII-output tests. — 707b1dc

### Phase 2: HTTP endpoint

- [x] 2.1 Add `/ascii-candles` handler, route, and failure/success tests. — adbd6c9
- [x] 2.2 Run full validation gate and review the implementation. — adbd6c9