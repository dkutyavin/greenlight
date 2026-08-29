package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"greenlight.dekutyavin.net/internal/app"
	"greenlight.dekutyavin.net/internal/config"
)

const version = "1.0.0"

func main() {
	var cfg config.Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()
	cfg.Version = version

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := app.Application{
		Config: cfg,
		Logger: logger,
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
