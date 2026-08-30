package app

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"greenlight.dekutyavin.net/internal/config"
)

func newTestApplication(t *testing.T) *Application {
	t.Helper()

	return &Application{
		Config: config.Config{
			Env: "testing",
		},
		Logger: slog.New(slog.DiscardHandler),
	}
}

func newTestApplicationWithLogs(t *testing.T) (*Application, *bytes.Buffer) {
	t.Helper()
	buf := &bytes.Buffer{}
	app := &Application{
		Config: config.Config{Env: "testing"},
		Logger: slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})),
	}

	return app, buf
}

func assertLogError(t *testing.T, buf *bytes.Buffer, errMsg string, method string, uri string) {
	t.Helper()

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
