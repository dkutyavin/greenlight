package app

import (
	"encoding/json"
	"errors"
	"log/slog"
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

	var entry struct {
		Level  string `json:"level"`
		Msg    string `json:"msg"`
		Method string `json:"method"`
		URI    string `json:"uri"`
	}

	err := json.Unmarshal(buf.Bytes(), &entry)
	if err != nil {
		t.Fatalf("log output is not valid json: %v\n%s", err, buf.String())
	}

	if entry.Level != slog.LevelError.String() {
		t.Errorf("want level: %q; got: %q", slog.LevelError.String(), entry.Level)
	}

	if entry.Msg != errMsg {
		t.Errorf("want msg: %q; got: %q", errMsg, entry.Msg)
	}

	if entry.Method != method {
		t.Errorf("want method: %q; got: %q", method, entry.Method)
	}

	if entry.URI != uri {
		t.Errorf("want URI: %q; got: %q", uri, entry.URI)
	}
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

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/any-path-will-do", nil)
	wantStatusCode := http.StatusInternalServerError
	message := envelope{"error": "the server encountered a problem and could not process your request"}
	wantBody, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("could not marshal wantBody: %v", err)
	}
	wantBody = append(wantBody, '\n')

	app.serverErrorResponse(rr, req, errors.New("some error"))

	// todo: check the logger calls right way
	if buf.Len() <= 0 {
		t.Errorf("there are no calls for logger")
	}

	if rr.Result().StatusCode != wantStatusCode {
		t.Errorf("want status: %d; got: %d", wantStatusCode, rr.Result().StatusCode)
	}

	gotBody := rr.Body.String()
	if gotBody != string(wantBody) {
		t.Errorf("want body: %q; got: %q", string(wantBody), gotBody)
	}
}
