# Variables
BINARY_NAME=tada-api
MAIN_PACKAGE=./cmd/api
BUILD_DIR=./bin
DOCKER_TAG=tada:latest

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.version=$(shell git describe --tags --always --dirty)"
BUILD_FLAGS=-v $(LDFLAGS)

.PHONY: all build clean test deps fmt vet lint docker-build docker-run help

## help: Display this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## all: Run all checks and build
all: fmt vet lint test build

## build: Build the application
build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

## clean: Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

## test: Run tests
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

## test-coverage: Run tests with coverage report
test-coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## deps: Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) verify

## fmt: Format code
fmt:
	$(GOCMD) fmt ./...

## vet: Run go vet
vet:
	$(GOCMD) vet ./...

## lint: Run golangci-lint
lint:
	golangci-lint run

## tidy: Tidy go modules
tidy:
	$(GOMOD) tidy

## run: Run the application
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

## watch: Run with hot reload using Air
watch:
	air

## docker-build: Build Docker image
docker-build:
	docker build -t $(DOCKER_TAG) .

## docker-run: Run Docker container
docker-run:
	docker run -p 8080:8080 $(DOCKER_TAG)

## docker-compose-up: Start services with docker-compose
docker-compose-up:
	docker-compose up -d

## docker-compose-down: Stop services
docker-compose-down:
	docker-compose down

## migration-up: Run database migrations
migration-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

## migration-down: Rollback database migrations
migration-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down

## seed: Seed the database
seed:
	$(GOCMD) run ./cmd/seed

## install-tools: Install development tools
install-tools:
	$(GOCMD) install github.com/air-verse/air@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	$(GOCMD) install github.com/swaggo/swag/cmd/swag@latest

## swagger: Generate swagger documentation
swagger:
	swag init -g cmd/api/main.go -o docs

## audit: Run security audit
audit:
	$(GOMOD) tidy -diff
	$(GOMOD) verify
	$(GOTEST) -race -vet=off ./...
	$(GOCMD) run golang.org/x/vuln/cmd/govulncheck@latest ./...

## ci: CI pipeline commands
ci: fmt vet lint test build

.DEFAULT_GOAL := help

