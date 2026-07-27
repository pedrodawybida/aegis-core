# Aegis Core - Build & Test Engine
PATH := /opt/homebrew/bin:/usr/local/bin:$(PATH)

.PHONY: help build test test-race run docker-build demo clean

help:
	@echo "🛡️  Aegis Core Management Commands:"
	@echo "  make build         - Build the aegis executable binary"
	@echo "  make test          - Run all unit and integration tests"
	@echo "  make test-race     - Run tests with race detector enabled"
	@echo "  make run           - Run Aegis Core locally"
	@echo "  make docker-build  - Build Docker container image"
	@echo "  make demo          - Execute interactive demonstration script"
	@echo "  make clean         - Clean built artifacts and logs"

build:
	go build -o bin/aegis cmd/aegis/main.go

test:
	go test -v ./...

test-race:
	go test -v -race ./...

run:
	go run cmd/aegis/main.go -config aegis.yaml -port 8080 -log audit_bacen.log

docker-build:
	docker build -t aegis-core:latest .

demo:
	./demo.sh

clean:
	rm -rf bin/ audit_bacen.log audit_test.log
