package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"greenlight.dekutyavin.net/internal/config"
)

func TestHealthCheckHandler(t *testing.T) {
	app := newTestApplication(t)
	testConfig := config.Config{
		Env:     "testing",
		Version: "1.0.0",
	}
	app.Config = testConfig

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)

	app.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want = %d", rr.Code, http.StatusOK)
	}

	gotContentType := rr.Header().Get("Content-Type")
	wantContentType := "application/json"

	if gotContentType != wantContentType {
		t.Errorf("want content type: %q; got content type:%q\n", wantContentType, gotContentType)
	}

	bodyBytes := rr.Body.Bytes()
	var resBody struct {
		Status     string `json:"status"`
		SystemInfo struct {
			Environment string `json:"environment"`
			Version     string `json:"version"`
		} `json:"system_info"`
	}

	err := json.Unmarshal(bodyBytes, &resBody)
	if err != nil {
		t.Fatalf("could not unmarshal json response body: %v", err)
	}

	if resBody.SystemInfo.Environment != testConfig.Env {
		t.Errorf("want environment: %q; got: %q", testConfig.Env, resBody.SystemInfo.Environment)
	}

	if resBody.Status != "available" {
		t.Errorf("want status: %q; got: %q", "available", resBody.Status)
	}

	if resBody.SystemInfo.Version != testConfig.Version {
		t.Errorf("want version: %q; got: %q", testConfig.Version, resBody.SystemInfo.Version)
	}
}
