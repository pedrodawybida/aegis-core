FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /aegis ./cmd/aegis/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /aegis .
COPY aegis.yaml .

EXPOSE 8080
CMD ["./aegis"]
