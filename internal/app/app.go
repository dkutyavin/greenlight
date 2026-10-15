package app

import (
	"log/slog"

	"greenlight.dekutyavin.net/internal/config"
	"greenlight.dekutyavin.net/internal/data"
)

type Application struct {
	Config config.Config
	Logger *slog.Logger
	Models data.Models
}
