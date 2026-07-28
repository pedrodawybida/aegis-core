# Nexo Hub - Build & Test Engine
PATH := /opt/homebrew/bin:/usr/local/bin:$(PATH)

.PHONY: help build test test-race run docker-build demo clean

help:
	@echo "🛡️  Nexo Hub Management Commands:"
	@echo "  make build         - Build the nexo executable binary"
	@echo "  make test          - Run all unit and integration tests"
	@echo "  make test-race     - Run tests with race detector enabled"
	@echo "  make run           - Run Nexo Hub locally"
	@echo "  make docker-build  - Build Docker container image"
	@echo "  make demo          - Execute interactive demonstration script"
	@echo "  make clean         - Clean built artifacts and logs"

build:
	go build -o bin/nexo cmd/nexo/main.go
	go build -o bin/nexo-init cmd/nexo-init/main.go

init:
	go run cmd/nexo-init/main.go

test:
	go test -v ./...

test-race:
	go test -v -race ./...

run:
	go run cmd/nexo/main.go -config nexo.yaml -port 8080 -log audit_bacen.log

docker-build:
	docker build -t nexo-hub:latest .

demo:
	./demo.sh

clean:
	rm -rf bin/ audit_bacen.log audit_test.log
