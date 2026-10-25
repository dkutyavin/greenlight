package app

import (
	"log/slog"

	"greenlight.dekutyavin.net/internal/config"
	"greenlight.dekutyavin.net/internal/data"
)

type Application struct {
	Config config.Config
	Logger *slog.Logger
	Models Models
}

type MovieModel interface {
	Insert(*data.Movie) error
	Get(int) (*data.Movie, error)
	Update(*data.Movie) error
	Delete(int) error
}

type Models struct {
	Movies MovieModel
}
