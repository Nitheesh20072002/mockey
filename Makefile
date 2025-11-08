# Helper Makefile (for *nix environments). On Windows use equivalent commands or WSL.
.PHONY: build run docker-up tidy

build:
	go build ./cmd/server

run:
	go run ./cmd/server

tidy:
	go mod tidy

docker-up:
	docker compose up --build
