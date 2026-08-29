package app

import (
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
