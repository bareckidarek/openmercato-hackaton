package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestKittyPickerUsesMinuteParityAndSelectedColour(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{name: "even minute", time: time.Date(2026, time.September, 19, 10, 2, 0, 0, time.UTC), want: "kitty-black-2.svg"},
		{name: "odd minute", time: time.Date(2026, time.September, 19, 10, 3, 0, 0, time.UTC), want: "kitty-red-1.svg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			picker := kittyPicker{
				now:  func() time.Time { return tt.time },
				pick: func(int) int { return 1 },
			}

			if got := picker.path(); !strings.HasSuffix(got, tt.want) {
				t.Fatalf("picker.path() = %q, want suffix %q", got, tt.want)
			}
		})
	}
}

func TestKittyPageRendersOneImage(t *testing.T) {
	app := &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		kitty: kittyPicker{
			now:  func() time.Time { return time.Date(2026, time.September, 19, 10, 2, 0, 0, time.UTC) },
			pick: func(int) int { return 0 },
		},
	}
	request := httptest.NewRequest(http.MethodGet, "/kitty", nil)
	response := httptest.NewRecorder()

	app.routes().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	body := response.Body.String()
	if got := strings.Count(body, "<img "); got != 1 {
		t.Fatalf("image count = %d, want 1", got)
	}
	if !strings.Contains(body, "/static/img/kitty-black-1.svg") {
		t.Fatalf("response does not contain selected kitty image: %s", body)
	}
}

func TestKittyPageRejectsUnknownPath(t *testing.T) {
	app := &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		kitty:  newKittyPicker(),
	}
	request := httptest.NewRequest(http.MethodGet, "/kitty/extra", nil)
	response := httptest.NewRecorder()

	app.routes().ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}
