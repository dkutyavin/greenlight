package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	app := application{
		config: config{env: "testing"},
		logger: slog.New(slog.DiscardHandler),
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)

	app.routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want = %d", rr.Code, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "environment: testing") {
		t.Errorf("body missing environment, got: %s", rr.Body.String())
	}
}
