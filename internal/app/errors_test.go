package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogError(t *testing.T) {
	t.Parallel()

	app, buf := newTestApplicationWithLogs(t)

	errMsg := "boom"
	method := http.MethodGet
	uri := "/v1/movies/42"
	req := httptest.NewRequest(method, uri, nil)

	app.logError(req, errors.New(errMsg))
	assertLogError(t, buf, errMsg, method, uri)
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	app := newTestApplication(t)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/any-path-will-do", nil)
	wantStatusCode := http.StatusBadRequest
	message := map[string]any{
		"key": "value",
	}
	wantBody, err := json.Marshal(envelope{"error": message})
	if err != nil {
		t.Fatalf("could not marshal wantBody: %v", err)
	}
	wantBody = append(wantBody, '\n')

	app.errorResponse(rr, req, wantStatusCode, message)

	if rr.Result().StatusCode != wantStatusCode {
		t.Errorf("want status: %d; got: %d", wantStatusCode, rr.Result().StatusCode)
	}

	gotBody := rr.Body.String()
	if gotBody != string(wantBody) {
		t.Errorf("want body: %q; got: %q", string(wantBody), gotBody)
	}
}

func TestServerErrorResponse(t *testing.T) {
	t.Parallel()

	app, buf := newTestApplicationWithLogs(t)

	method := http.MethodGet
	uri := "/v1/any-path-will-do"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(method, uri, nil)
	wantStatusCode := http.StatusInternalServerError
	errMsg := "some error message"
	message := envelope{"error": "the server encountered a problem and could not process your request"}
	wantBody, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("could not marshal wantBody: %v", err)
	}
	wantBody = append(wantBody, '\n')

	app.serverErrorResponse(rr, req, errors.New(errMsg))

	assertLogError(t, buf, errMsg, method, uri)
	if rr.Result().StatusCode != wantStatusCode {
		t.Errorf("want status: %d; got: %d", wantStatusCode, rr.Result().StatusCode)
	}

	gotBody := rr.Body.String()
	if gotBody != string(wantBody) {
		t.Errorf("want body: %q; got: %q", string(wantBody), gotBody)
	}
}
