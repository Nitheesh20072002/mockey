# Exam Platform API (Gin)

This repository is a starter scaffold for an exam/test platform backend using Go + Gin.

Goals
- Provide REST endpoints for Auth, Admin (exams/questions), Test lifecycle, and Analytics.
- Use Postgres as the primary DB, keep DB layer extensible.
- Use gvm to manage Go versions.
- Be containerizable (Docker + docker-compose).

What's included in this scaffold
- `cmd/server/main.go` — minimal Gin server with health check
- `internal/server/router.go` — router setup
- `go.mod` — go module file with Gin dependency
- `Dockerfile` & `docker-compose.yml` — run app + Postgres
- `.env.example` — example env variables
- `.go-version` — recommended Go version string for gvm
- `Makefile` — helper targets

First steps (recommended)
1. Install gvm (follow project doc for your platform), then install the recommended Go version:

   # Example (Linux/macOS gvm syntax)
   gvm install go1.20.7
   gvm use go1.20.7 --default

   For Windows, use gvm-windows or manage Go versions with another tool if you prefer.

2. Copy `.env.example` to `.env` and set values.

3. Locally (if you have Go installed):

   go mod tidy
   go run ./cmd/server

4. Using Docker (recommended):

   docker compose up --build

This brings up the API on port 8080 and Postgres on 5432.

Next: review the implementation plan in `IMPLEMENTATION_PLAN.md` (to be added) and confirm to proceed with endpoint implementations, DB models, migrations, and auth.
