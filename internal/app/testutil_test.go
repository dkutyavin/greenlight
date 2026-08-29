package app

import (
	"log/slog"
	"testing"

	"greenlight.dekutyavin.net/internal/data"
)

func newTestApplication(t *testing.T) *Application {
	t.Helper()

	return &Application{
		Config: data.Config{
			Env: "testing",
		},
		Logger: slog.New(slog.DiscardHandler),
	}
}
