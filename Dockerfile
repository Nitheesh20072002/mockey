# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app

# Install git for `go get` modules if needed
RUN apk add --no-cache git ca-certificates

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /exam-api ./cmd/server

# Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /exam-api /exam-api
EXPOSE 8080
ENTRYPOINT ["/exam-api"]
