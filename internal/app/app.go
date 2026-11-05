package app

import (
	"log/slog"
	"sync"

	"greenlight.dekutyavin.net/internal/config"
	"greenlight.dekutyavin.net/internal/data"
	"greenlight.dekutyavin.net/internal/mailer"
)

type Application struct {
	Config config.Config
	Logger *slog.Logger
	Models Models
	Mailer *mailer.Mailer
	WG     sync.WaitGroup
}

type MovieModel interface {
	Insert(*data.Movie) error
	Get(int) (*data.Movie, error)
	Update(*data.Movie) error
	Delete(int) error
	GetAll(string, []string, data.Filters) ([]*data.Movie, data.Metadata, error)
}

type UserModel interface {
	Insert(*data.User) error
	GetByEmail(string) (*data.User, error)
	Update(*data.User) error
}

type Models struct {
	Movies MovieModel
	Users  UserModel
}
