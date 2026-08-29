package app

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"greenlight.dekutyavin.net/internal/data"
)

func TestHealthCheckHandler(t *testing.T) {
	app := Application{
		Config: data.Config{Env: "testing"},
		Logger: slog.New(slog.DiscardHandler),
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)

	app.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want = %d", rr.Code, http.StatusOK)
	}

	gotContentType := rr.Header().Get("Content-Type")
	wantContentType := "application/json"

	if gotContentType != wantContentType {
		t.Errorf("want content type: %s; got content type:%s\n", wantContentType, gotContentType)
	}

	bodyBytes := rr.Body.Bytes()
	var actual data.HealthCheckResponse

	err := json.Unmarshal(bodyBytes, &actual)
	if err != nil {
		t.Fatalf("could not unmarshal json response body: %v", err)
	}

	fmt.Printf("%v", actual)

	if actual.Environment != app.Config.Env {
		t.Errorf("want environment: %q; got: %q", app.Config.Env, actual.Environment)
	}

	if actual.Status != "available" {
		t.Errorf("want status: %q; got: %q", "available", actual.Status)
	}

	if actual.Version != app.Config.Version {
		t.Errorf("want version: %q; got: %q", app.Config.Version, actual.Version)
	}
}
