package main

import "net/http"

func (app *application) asciiCandlesPage(w http.ResponseWriter, r *http.Request) {
	candles, err := app.market.Daily(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(renderCandles(candles)))
}