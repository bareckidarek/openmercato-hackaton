package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeMarketClient struct {
	candles []ohlc
	err     error
}

func (f fakeMarketClient) Daily(context.Context) ([]ohlc, error) {
	return f.candles, f.err
}

func TestASCIICandlesPage(t *testing.T) {
	app := &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		market: fakeMarketClient{candles: []ohlc{{Date: time.Now(), Open: 1, High: 2, Low: 0, Close: 2}}},
	}
	request := httptest.NewRequest(http.MethodGet, "/ascii-candles", nil)
	response := httptest.NewRecorder()
	app.routes().ServeHTTP(response, request)
	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("status/content type = %d/%q", response.Code, response.Header().Get("Content-Type"))
	}
	if !strings.Contains(response.Body.String(), "SPY - recent daily candles") {
		t.Fatalf("unexpected body: %s", response.Body.String())
	}
}

func TestASCIICandlesPageReturnsServerError(t *testing.T) {
	app := &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		market: fakeMarketClient{err: io.ErrUnexpectedEOF},
	}
	request := httptest.NewRequest(http.MethodGet, "/ascii-candles", nil)
	response := httptest.NewRecorder()
	app.routes().ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}

func TestYahooClientDailyValidatesAndLimitsData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"chart":{"result":[{"timestamp":[1,2],"indicators":{"quote":[{"open":[1,2],"high":[3,4],"low":[0,1],"close":[2,3]}]}}]}}`)
	}))
	defer server.Close()
	got, err := (yahooClient{httpClient: server.Client(), url: server.URL}).Daily(context.Background())
	if err != nil || len(got) != 2 || got[0].Close != 2 {
		t.Fatalf("Daily() = %#v, %v", got, err)
	}
}

func TestYahooClientRejectsInvalidData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"chart":{"result":[{"timestamp":[1],"indicators":{"quote":[{"open":[4],"high":[3],"low":[1],"close":[2]}]}}]}}`)
	}))
	defer server.Close()
	if _, err := (yahooClient{httpClient: server.Client(), url: server.URL}).Daily(context.Background()); err == nil {
		t.Fatal("Daily() accepted invalid OHLC data")
	}
}

func TestRenderCandlesIsASCIIAndShowsTrend(t *testing.T) {
	output := renderCandles([]ohlc{
		{Date: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), Open: 10, High: 12, Low: 9, Close: 11},
		{Date: time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC), Open: 11, High: 12, Low: 8, Close: 9},
	})
	if !strings.Contains(output, "#") || !strings.Contains(output, ":") || !strings.Contains(output, "09/01") {
		t.Fatalf("renderCandles() missing expected chart content:\n%s", output)
	}
	for i := 0; i < len(output); i++ {
		if output[i] > 127 {
			t.Fatalf("renderCandles() emitted non-ASCII byte %d", output[i])
		}
	}
}