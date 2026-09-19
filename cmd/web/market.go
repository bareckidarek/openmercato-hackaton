package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

const yahooSPYURL = "https://query1.finance.yahoo.com/v8/finance/chart/SPY?range=1mo&interval=1d"

type ohlc struct {
	Date       time.Time
	Open, High float64
	Low, Close float64
}

type marketDataClient interface {
	Daily(context.Context) ([]ohlc, error)
}

type yahooClient struct {
	httpClient *http.Client
	url        string
}

type yahooChartResponse struct {
	Chart struct {
		Result []struct {
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open  []*float64 `json:"open"`
					High  []*float64 `json:"high"`
					Low   []*float64 `json:"low"`
					Close []*float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error json.RawMessage `json:"error"`
	} `json:"chart"`
}

func (c yahooClient) Daily(ctx context.Context) ([]ohlc, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "openmercato-ascii-candles/1.0")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("market data returned status %d", resp.StatusCode)
	}

	var payload yahooChartResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload); err != nil {
		return nil, err
	}
	if len(payload.Chart.Result) != 1 {
		return nil, errors.New("market data contained no chart")
	}
	result := payload.Chart.Result[0]
	if len(result.Indicators.Quote) != 1 {
		return nil, errors.New("market data contained no quotes")
	}
	quote := result.Indicators.Quote[0]
	n := len(result.Timestamp)
	if n == 0 || len(quote.Open) != n || len(quote.High) != n || len(quote.Low) != n || len(quote.Close) != n {
		return nil, errors.New("market data arrays have inconsistent lengths")
	}
	out := make([]ohlc, 0, n)
	for i, timestamp := range result.Timestamp {
		if quote.Open[i] == nil || quote.High[i] == nil || quote.Low[i] == nil || quote.Close[i] == nil {
			return nil, errors.New("market data contains null values")
		}
		candle := ohlc{
			Date:  time.Unix(timestamp, 0).UTC(),
			Open:  *quote.Open[i],
			High:  *quote.High[i],
			Low:   *quote.Low[i],
			Close: *quote.Close[i],
		}
		if !validOHLC(candle) {
			return nil, errors.New("market data contains invalid prices")
		}
		out = append(out, candle)
	}
	if len(out) > 30 {
		out = out[len(out)-30:]
	}
	return out, nil
}

func validOHLC(c ohlc) bool {
	for _, value := range []float64{c.Open, c.High, c.Low, c.Close} {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return false
		}
	}
	return c.Low <= c.Open && c.Low <= c.Close && c.High >= c.Open && c.High >= c.Close && c.Low <= c.High
}

func renderCandles(candles []ohlc) string {
	const rows = 12
	var b strings.Builder
	b.WriteString("SPY - recent daily candles (USD)\n")
	b.WriteString("Source: Yahoo Finance chart data | Not investment advice\n")
	b.WriteString("Legend: # bullish body  : bearish body  | wick\n")
	if len(candles) == 0 {
		b.WriteString("No market data available.\n")
		return b.String()
	}
	minimum, maximum := candles[0].Low, candles[0].High
	for _, candle := range candles[1:] {
		minimum = math.Min(minimum, candle.Low)
		maximum = math.Max(maximum, candle.High)
	}
	if maximum == minimum {
		maximum = minimum + 1
	}
	scale := func(value float64) int {
		return rows - 1 - int(math.Round((value-minimum)/(maximum-minimum)*float64(rows-1)))
	}
	grid := make([][]byte, rows)
	for row := range grid {
		grid[row] = []byte(strings.Repeat(" ", len(candles)*3))
	}
	for i, candle := range candles {
		high, low := scale(candle.High), scale(candle.Low)
		bodyTop, bodyBottom := scale(math.Max(candle.Open, candle.Close)), scale(math.Min(candle.Open, candle.Close))
		if bodyTop == bodyBottom {
			bodyBottom = min(rows-1, bodyBottom+1)
		}
		marker := byte(':')
		if candle.Close >= candle.Open {
			marker = '#'
		}
		for row := low; row <= high; row++ {
			grid[row][i*3+1] = '|'
		}
		for row := bodyTop; row <= bodyBottom; row++ {
			grid[row][i*3] = marker
			grid[row][i*3+1] = marker
			grid[row][i*3+2] = marker
		}
	}
	for row, line := range grid {
		b.WriteString(fmt.Sprintf("%7.0f ", maximum-(maximum-minimum)*float64(row)/float64(rows-1)))
		b.Write(line)
		b.WriteByte('\n')
	}
	b.WriteString("Dates: ")
	for i, candle := range candles {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(candle.Date.Format("01/02"))
	}
	b.WriteByte('\n')
	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}