# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /nexo ./cmd/nexo/main.go

# Production zero-overhead scratch runtime stage
FROM scratch

WORKDIR /app
COPY --from=builder /nexo /app/nexo
COPY nexo.yaml /app/nexo.yaml

EXPOSE 8080
ENTRYPOINT ["/app/nexo"]
