# Helper Makefile (for *nix environments). On Windows use equivalent commands or WSL.
.PHONY: build run docker-up tidy

build:
	go build ./cmd/server

run:
	go run ./cmd/server

tidy:
	go mod tidy

dev:
	docker compose up --build
