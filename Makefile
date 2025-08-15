APP_NAME := pick
SRC := cmd/main.go
OUTPUT := $(APP_NAME)

all: build

build: $(SRC)
	@echo "Building the project..."
	go build -o $(OUTPUT) $(SRC)

run: 
	@echo "Running the application..."
	CONFIG_PATH=./config/config.toml go run $(SRC)

debug:
	@echo "Running the application with -race"
	go run -race $(SRC)

swag:
	@echo "Generating swagger docs..."
	swag init -g ./cmd/main.go -o ./docs

lint:
	@echo "Running golangci-lint locally..."
	golangci-lint run --timeout=5m

lint-fix:
	@echo "Running golangci-lint with fix..."
	golangci-lint run --timeout=5m --fix

lint-ci:
	@echo "Running golangci-lint in CI mode..."
	golangci-lint run --timeout=5m --out-format=github-actions

test:
	@echo "Running tests with race detector and coverage report..."
	go test ./internal/... -race -coverprofile=coverage.out -covermode=atomic
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

test-e2e:
	@echo "Running e2e test"
	go test -v ./tests/e2e/...

test-integration:
	@echo "Running integration test"
	go test -v ./tests/integration/...

help:
	@echo "Makefile targets:"
	@echo "  build   			- Build the Go project"
	@echo "  run     			- Build and run the project"
	@echo "  debug   			- Run the application with debug mode"
	@echo "  swag   			- Generate Swagger Template"
	@echo "  test    			- Run tests with race detector and generate coverage report"
	@echo "  test-e2e   		- Run end-to-end Go tests"
	@echo "  test-integration 	- Run integration tests"
	@echo "  lint         		- Run golangci-lint locally"
	@echo "  lint-fix     		- Run golangci-lint and automatically fix issues"
	@echo "  lint-ci      		- Run golangci-lint in CI mode for GitHub Actions"
	@echo "  help    			- Show this help message"