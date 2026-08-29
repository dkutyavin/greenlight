package app

import (
	"log/slog"

	"greenlight.dekutyavin.net/internal/data"
)

type Application struct {
	Config data.Config
	Logger *slog.Logger
}
