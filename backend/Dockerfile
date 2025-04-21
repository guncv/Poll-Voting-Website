# Dockerfile for Backend

FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./cmd/main.go

# Final image
FROM alpine:3.18

WORKDIR /app

# Install necessary dependencies
RUN apk add --no-cache ca-certificates

# Copy the binary and environment variables
COPY --from=builder /app/app /app/app
COPY .env /app/.env

EXPOSE 8080
ENTRYPOINT ["/app/app"]
