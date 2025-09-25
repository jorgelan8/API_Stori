# API Stori - Makefile

.PHONY: help build run stop logs test clean docker-build docker-run docker-stop

# Default target
help:
	@echo "API Stori - Available commands:"
	@echo ""
	@echo "🐳 Docker commands:"
	@echo "  make start     - Start API with Docker Compose"
	@echo "  make stop      - Stop API and remove containers"
	@echo "  make logs      - Show API logs"
	@echo "  make restart   - Restart API"
	@echo "  make build     - Build Docker image"
	@echo ""
	@echo "🔧 Development commands:"
	@echo "  make dev       - Run API locally (requires Go)"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"
	@echo ""
	@echo "🧪 Testing commands:"
	@echo "  make test-unit        - Run unit tests with coverage"
	@echo "  make test-integration - Run integration tests with coverage"
	@echo "  make test-load        - Run load tests with coverage"
	@echo "  make test-performance - Run performance tests with coverage"
	@echo "  make test-all         - Run all tests with coverage"
	@echo "  make coverage-report  - Generate HTML coverage report"
	@echo "  make coverage-summary - Show coverage summary"
	@echo "  make test-api         - Test API endpoints (Docker)"
	@echo "  make test-csv         - Test CSV migration (Docker)"
	@echo "  make test-balance     - Test balance endpoint (Docker)"

# Docker commands
start:
	@echo "🚀 Starting API Stori with Docker Compose..."
	@./start.sh

stop:
	@echo "🛑 Stopping API Stori..."
	@./stop.sh

restart: stop start

logs:
	@echo "📋 Showing API logs..."
	@docker-compose logs -f

build:
	@echo "📦 Building Docker image..."
	@docker-compose build

# Development commands
dev:
	@echo "🔧 Running API locally..."
	@go run cmd/api/main.go

test:
	@echo "🧪 Running tests..."
	@go test ./...

clean:
	@echo "🧹 Cleaning build artifacts..."
	@go clean
	@docker-compose down --rmi all --volumes --remove-orphans

# Testing commands
test-unit:
	@echo "🔬 Running unit tests with coverage..."
	@go test -v ./internal/services/... ./internal/handlers/... -coverprofile=coverage_unit.out

test-integration:
	@echo "🔗 Running integration tests with coverage..."
	@go test -v ./tests/integration/... -cover -coverpkg=./... -coverprofile=coverage_integration.out

test-all:
	@echo "🧪 Running all tests with coverage..."
	@./run_tests.sh

coverage-report:
	@echo "📊 Generating coverage report..."
	@go tool cover -html=coverage_all.out -o coverage_report.html
	@echo "📈 Coverage report generated: coverage_report.html"

coverage-summary:
	@echo "📊 Coverage Summary:"
	@echo "Unit Tests:"
	@go tool cover -func=coverage_unit.out | tail -1
	@echo "Integration Tests:"
	@go tool cover -func=coverage_integration.out | tail -1
	@echo "All Tests:"
	@go tool cover -func=coverage_all.out | tail -1

test-load:
	@echo "🧪 Running load tests with coverage..."
	@go test -v ./tests/load/... -cover -coverpkg=./... -coverprofile=coverage_load.out

test-performance:
	@echo "🚀 Running performance tests with coverage..."
	@go test -v ./tests/performance/... -cover -coverpkg=./... -coverprofile=coverage_performance.out

test-api:
	@echo "🧪 Testing API endpoints (Docker on port 8081)..."
	@echo "Testing health endpoint..."
	@curl -s http://localhost:8081/api/v1/health | jq .
	@echo ""
	@echo "Testing root endpoint..."
	@curl -s http://localhost:8081/ | jq .

test-csv:
	@echo "🧪 Testing CSV migration (Docker on port 8081)..."
	@curl -X POST http://localhost:8081/api/v1/migrate \
		-F "csv_file=@examples/sample_transactions.csv" | jq .

test-balance:
	@echo "🧪 Testing balance endpoint (Docker on port 8081)..."
	@curl -s "http://localhost:8081/api/v1/users/1001/balance" | jq .

# Docker run commands
docker-run:
	@echo "🐳 Running API with Docker..."
	@docker run -p 8081:8080 jps-api-stori

docker-stop:
	@echo "🛑 Stopping Docker containers..."
	@docker stop $$(docker ps -q --filter ancestor=jps-api-stori) 2>/dev/null || true
