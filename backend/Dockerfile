# syntax=docker/dockerfile:1

# 1) Builder Stage: use ARM64 Go so it runs natively
FROM --platform=linux/arm64 golang:1.24-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Cross-compile to AMD64 (static binary)
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./cmd/main.go

# 2) Runtime Stage: pull the AMD64 Alpine image
FROM --platform=linux/amd64 alpine:3.18

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app /app/app
COPY .env            /app/.env

EXPOSE 8080
ENTRYPOINT ["/app/app"]
