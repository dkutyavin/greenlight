package app

import (
	"log/slog"

	"greenlight.dekutyavin.net/internal/config"
)

type Application struct {
	Config config.Config
	Logger *slog.Logger
}
