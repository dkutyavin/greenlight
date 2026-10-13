ifneq (,$(wildcard .env))
include .env
export
endif

.PHONY: build run test db-up db-down prepare-env

db-up:
	docker compose up -d --wait

db-down:
	docker compose down

prepare-env:
	cp .env.example .env

build:
	go build -o=./bin/api ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...
