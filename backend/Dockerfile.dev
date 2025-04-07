FROM golang:1.23-alpine

# Install system dependencies
RUN apk add --no-cache git bash

# Set working directory
WORKDIR /usr/src/app

# Install development tools: Air for hot-reloading and swag for Swagger docs.
RUN go install github.com/cosmtrek/air@v1.41.0 && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Ensure Go binaries are in PATH
ENV PATH="/go/bin:${PATH}"

# Copy go.mod and go.sum first to leverage Docker caching for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Optionally generate Swagger docs (uncomment if needed)
# RUN swag init

# Tidy up Go modules
RUN go mod tidy

EXPOSE 8080

# Start Air with the specified configuration file for hot-reloading
CMD ["air", "-c", ".air.toml"]
