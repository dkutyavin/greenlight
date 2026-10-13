package config

import "time"

type Config struct {
	Port    int
	Env     string
	Version string
	DB      struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
}
