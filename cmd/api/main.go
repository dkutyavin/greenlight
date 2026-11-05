package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"strconv"
	"time"

	"greenlight.dekutyavin.net/internal/app"
	"greenlight.dekutyavin.net/internal/config"
	"greenlight.dekutyavin.net/internal/data"
	"greenlight.dekutyavin.net/internal/mailer"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

func main() {
	var cfg config.Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.DB.DSN, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.SMTP.Host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.SMTP.Port, "smtp-port", envInt("GREENLIGHT_SMTP_PORT", 25), "SMTP port")
	flag.StringVar(&cfg.SMTP.Username, "smtp-username", os.Getenv("GREENLIGHT_SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.SMTP.Password, "smtp-password", os.Getenv("GREENLIGHT_SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.SMTP.Sender, "smtp-sender", "Greenlight <no-reply@greenlight.dekutyavin.net>", "SMTP sender")

	flag.Parse()
	cfg.Version = version

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	mailer, err := mailer.New(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.Sender,
	)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := app.Application{
		Config: cfg,
		Logger: logger,
		Models: app.Models{
			Movies: data.NewMoviesModel(db),
			Users:  data.NewUserModel(db),
		},
		Mailer: mailer,
	}

	err = app.Serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.DB.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// envInt reads the environment variable named key as an integer, falling back
// to fallback when the variable is unset or not a valid integer.
func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return n
}
