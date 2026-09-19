# Add `/ascii-candles` S&P 500 ETF chart endpoint

## 📝 TLDR
Future behavior: add a `GET /ascii-candles` endpoint that fetches recent daily OHLC data for the SPDR S&P 500 ETF Trust (`SPY`) and renders a compact candlestick chart as plain ASCII text. The endpoint is intended for quick terminal-friendly inspection, not trading advice or a replacement for a market-data terminal.

## 📝 Problem Statement
The application currently has no market-data endpoint. A small, dependency-light text endpoint can make recent SPY price movement observable without adding a charting frontend. Yahoo Finance's chart endpoint identifies SPY as an ETF and supplies daily open, high, low, close, volume, and timestamps suitable for this display.

## 📝 Proposed Solution
Add an HTTP handler that requests a bounded recent daily window from Yahoo Finance's chart API, validates the response, maps each trading day to one fixed-width ASCII candle, and returns `text/plain; charset=utf-8`. Up candles use `#`, down candles use `:`, and the high/low wick uses `|`; a legend and date/close labels make the output interpretable.

Alternatives rejected: an external charting library adds unnecessary browser/runtime weight; a persistent market-data store is out of scope; scraping HTML is less stable than Yahoo's JSON chart endpoint.

## 📝 Architecture
`GET /ascii-candles` uses a small injected market-data client and renderer in `cmd/web`, following the existing route, middleware, and error-response pipeline. The client uses Go's standard HTTP transport with a timeout; the renderer is pure and independently testable. No database, migration, background job, or frontend asset is required.

The upstream request should be made at request time with a bounded range (default: 1 month, daily interval). The handler must not expose arbitrary upstream URLs or query parameters.

## 📝 Data Model
No persisted data or schema changes. The internal response model contains timestamp/date, open, high, low, close, and optional volume. Upstream payloads are transient and must not be logged in full.

## 📝 API Contracts
### `GET /ascii-candles`

- Query parameters: none in the initial version.
- Success: `200 OK`, `Content-Type: text/plain; charset=utf-8`.
- Body: title identifying `SPY`, source/disclaimer line, legend, one aligned candle column per returned trading day, and a date/close footer or label row.
- The output contains at most 30 daily candles and preserves chronological order.
- A candle is bullish when `close >= open` and bearish otherwise. The body and wick are scaled against the returned window's high/low range; a flat candle remains visible with a minimum body height.
- Upstream data must be rejected when required arrays are absent, lengths differ, values are non-finite, or `high < max(open, close)` / `low > min(open, close)`.
- Upstream timeout, transport failure, non-2xx response, invalid JSON, or invalid market data returns the repository-standard `500` response without leaking upstream details.
- The endpoint is informational only and must include a concise “not investment advice” disclaimer.

Example shape:

```text
SPY — recent daily candles (USD)
Source: Yahoo Finance chart data | Not investment advice
Legend: # bullish body  : bearish body  | wick
  774          |  |
  768       |  ##|
  762    |  :##:|
  756       |  ::|
       09/01 09/02 09/03 ...
```

## 📝 UI/UX
This is a plain-text endpoint rather than a browser UI. Preserve whitespace with a monospace response, keep line width bounded for terminal use, and make the title, source, legend, and disclaimer readable without color. Empty or unavailable data must produce an explicit error response rather than an empty chart.

## 📝 Edge Cases & Failure Scenarios
- Weekend/holiday windows may contain fewer than the maximum number of candles; render the available chronological observations.
- A single trading-day window uses a stable non-zero vertical scale.
- Equal open/close values render a visible one-row body.
- Extreme prices or decimals are normalized to a bounded display precision without changing the source values used for scaling.
- Yahoo rate limiting, schema drift, or malformed values returns a safe server error and logs only a concise operational error.
- A slow upstream request is cancelled by the client or timeout; no goroutine or response write continues afterward.

## 📝 Risks & Impact Review
The endpoint adds a runtime dependency on an external market-data service, so availability and freshness are not guaranteed. The contract is deliberately narrow and read-only; no trading, authentication, financial advice, or persisted user data is introduced. Tests must use a fake HTTP client/server and must not call Yahoo Finance. The source attribution and disclaimer reduce, but do not eliminate, interpretation risk.

## Resolved assumptions (autonomous defaults)

| Question | Default | Rationale |
|---|---|---|
| Which instrument? | SPY, the SPDR S&P 500 ETF Trust | The brief says S&P 500 ETF index; SPY is the least ambiguous widely traded ETF representation. |
| Which data window and interval? | Most recent one month of daily candles, capped at 30 observations | Keeps output compact and limits upstream work while showing a useful pattern. |
| Should callers select symbols, ranges, or intervals? | No; defer query customization | Smallest public contract and lowest misuse/blast radius for the first version. |
| Should data be cached or persisted? | No; request-time fetch only | Avoids stale-data policy and schema/storage changes. |
| Should the endpoint provide trading signals? | No; display descriptive candles only | Avoids financial-advice scope and keeps the feature reversible. |

## 📋 Phasing

### Phase 1: Data contract and renderer
Implement the validated internal OHLC model, Yahoo response decoder, bounded HTTP client, ASCII candle renderer, and unit tests with fake upstream responses.

### Phase 2: HTTP endpoint and resilience
Register `/ascii-candles`, wire the handler through existing middleware, cover success and upstream failure behavior, and run the full validation gate.

## 📋 Implementation Plan

1. **Add market-data types and client** — define the minimal OHLC response model, decode/validate the Yahoo chart payload, enforce timeout and response-size limits, and test malformed/partial/upstream failure cases with a fake server.
2. **Add ASCII candle renderer** — implement deterministic scaling, bullish/bearish bodies, wicks, labels, legend, disclaimer, width limits, and tests for flat, single-day, mixed-trend, and empty data.
3. **Add `/ascii-candles` handler and route** — wire the injected client and renderer into the existing application, return plain text on success, use the standard server-error path on failure, and add HTTP route tests.
4. **Run verification** — run `make test`, `make build`, `git diff --check`, and manually re-read the endpoint contract and compatibility surfaces.

## 🔗 Research sources
- Yahoo Finance chart API response for `SPY`, daily interval, recent one-month range: `https://query1.finance.yahoo.com/v8/finance/chart/SPY?range=1mo&interval=1d` (retrieved 2026-09-19). It identifies the instrument as “State Street SPDR S&P 500 ETF Trust” and exposes daily OHLC arrays.
- The response observed during research included 22 daily observations, `currency: USD`, and `instrumentType: ETF`; values are volatile and time-dependent, so no sample quote is treated as a product requirement.