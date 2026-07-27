# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /aegis ./cmd/aegis/main.go

# Production zero-overhead scratch runtime stage
FROM scratch

WORKDIR /app
COPY --from=builder /aegis /app/aegis
COPY aegis.yaml /app/aegis.yaml

EXPOSE 8080
ENTRYPOINT ["/app/aegis"]
