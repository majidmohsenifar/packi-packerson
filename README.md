# 📦 Pack Service

A small test project that provides two HTTP endpoints to calculate "packs" for a given ordered quantity.

- **`/correct`** – Optimized pack calculation in Go
- **`/incorrect`** – Naive pack calculation in Go
- **Unit Tests** – Written in Go
- **Integration Tests** – Written in Go
- **E2E Tests** – Written in Go


## 🚀 Run the Project

```bash
CONFIG_PATH=./config/config.toml make run
```

Demo Swagger
```bash
http://localhost:8080/swagger/index.html
```

The server will start on the port defined in `config.toml`.
```bash
curl "http://localhost:8080/correct?x=12001"

curl "http://localhost:8080/incorrect?x=12001"
```


## 🧪 Tests

This project has **three levels of automated tests**:

### 1️⃣ Unit Tests (Go)  
Test the core pack calculation logic **in isolation**, without HTTP or other dependencies.

**Files:**  
- internal/pack/pack_test.go  

**Run:**  
```bash
make test
```

### 2️⃣ Integration Tests (Go)  
Test the HTTP API by sending requests directly to the server (requires the API to be running).

**Files:**  
- tests/integration/api_integration_test.go

**Run:**  
```bash
make test-integration
```

### 3️⃣ End-to-End (E2E)
Simulate **real user requests** over HTTP, validating API responses against expected outputs.

**Files:**  
- tests/e2e/e2e_test.go

**Run:**  
```bash
make test-e2e
```

## 🐳 Docker

Build and run with Docker:
```bash
docker build -t pack-service .
docker run -p 8080:8080 pack-service
```


## 📂 Project Structure

```
.
├── internal/pack                   # Go business logic (unit tests here)
├── api                             # HTTP API handlers
├── tests
│   ├── integration                 # directory for go integration tests
│   └── e2e                         # directory for E2E tests
├── config                          # Configuration files
├── shell.nix                       # Nix shell with Go, go-swag, golangci
├── .golangci.yml                   # Golang CI config
├── Dockerfile
├── go.mod / go.sum
└── Makefile
```


## ⚙️ Requirements

- Go ≥ 1.24  
- **Nix** (optional, to get all tools in one shell)


## ❄️ Nix Environment

This project includes a shell.nix file that sets up **Go, go-swag, golangci ** in a reproducible development environment.

```bash
nix-shell
# Inside Nix shell, you can run:
# Run unit tests
make test

# Run integration tests
make test-integration

# Run E2E tests
make test-e2e

```