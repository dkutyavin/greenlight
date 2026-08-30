package app

import (
	"bytes"
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
