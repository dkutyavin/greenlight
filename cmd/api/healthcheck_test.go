package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

	gotContentType := rr.Header().Get("Content-Type")
	wantContentType := "application/json"

	if gotContentType != wantContentType {
		t.Errorf("want content type: %s; got content type:%s\n", wantContentType, gotContentType)
	}

	bodyBytes := rr.Body.Bytes()
	var actual HealthCheckResponse

	err := json.Unmarshal(bodyBytes, &actual)
	if err != nil {
		t.Fatalf("could not unmarshal json response body: %v", err)
	}

	fmt.Printf("%v", actual)

	if actual.Environment != app.config.env {
		t.Errorf("want environment: %q; got: %q", app.config.env, actual.Environment)
	}

	if actual.Status != "available" {
		t.Errorf("want status: %q; got: %q", "available", actual.Status)
	}

	if actual.Version != version {
		t.Errorf("want version: %q; got: %q", version, actual.Version)
	}
}
