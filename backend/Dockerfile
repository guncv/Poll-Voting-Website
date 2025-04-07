FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/main.go

# Final image
FROM scratch
COPY --from=builder /app/app /app/app

EXPOSE 8080
ENTRYPOINT ["/app/app"]
