ifneq (,$(wildcard .env))
include .env
export
endif

.PHONY: build run test db-up db-down prepare-env migrate-up migrate-down

DB_DSN ?= postgres://${DB_USER}:${DB_PASS}@localhost:5432/${DB_NAME}?sslmode=disable

db-up:
	docker compose up -d --wait

db-down:
	docker compose down

prepare-env:
	cp .env.example .env

build:
	go build -o=./bin/api ./cmd/api

run:
	GREENLIGHT_DB_DSN='${DB_DSN}' go run ./cmd/api

test:
	go test ./...

migrate-up:
	migrate -path=./migrations -database='${DB_DSN}' up

migrate-down:
	migrate -path=./migrations -database='${DB_DSN}' down
