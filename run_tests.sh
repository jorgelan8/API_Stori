#!/bin/bash

# API Stori - Test Runner Script
echo "🧪 Running API Stori Tests..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    print_error "Go version $GO_VERSION is not supported. Please install Go $REQUIRED_VERSION or higher."
    exit 1
fi

print_status "Go version $GO_VERSION detected"

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy
if [ $? -ne 0 ]; then
    print_error "Failed to install dependencies"
    exit 1
fi

print_status "Dependencies installed"

# Run unit tests
echo ""
echo "🔬 Running unit tests with coverage..."
go test -v ./internal/services/... ./internal/handlers/... -coverprofile=coverage_unit.out
if [ $? -ne 0 ]; then
    print_error "Unit tests failed"
    exit 1
fi

print_status "Unit tests passed"

# Run integration tests
echo ""
echo "🔗 Running integration tests with coverage..."
go test -v ./tests/integration/... -cover -coverpkg=./... -coverprofile=coverage_integration.out
if [ $? -ne 0 ]; then
    print_error "Integration tests failed"
    exit 1
fi

print_status "Integration tests passed"

# Run all tests with coverage
echo ""
echo "📊 Running all tests with coverage..."
go test -v ./... -cover -coverpkg=./... -coverprofile=coverage_all.out
if [ $? -ne 0 ]; then
    print_error "Some tests failed"
    exit 1
fi

# Generate coverage report
echo ""
echo "📈 Generating coverage report..."
go tool cover -html=coverage_all.out -o coverage_report.html
if [ $? -eq 0 ]; then
    print_status "Coverage report generated: coverage_report.html"
else
    print_warning "Failed to generate HTML coverage report"
fi

# Show coverage summary
echo ""
echo "📊 Coverage Summary:"
go tool cover -func=coverage_all.out | tail -1

# Clean up coverage files
echo ""
echo "🧹 Cleaning up..."
rm -f coverage_unit.out coverage_integration.out coverage_all.out

print_status "All tests completed successfully!"
echo ""
echo "🎉 Test Summary:"
echo "   ✅ Unit tests: PASSED"
echo "   ✅ Integration tests: PASSED"
echo "   ✅ Coverage report: coverage_report.html"
echo ""
echo "📋 To run specific tests:"
echo "   go test -v ./internal/services/... ./internal/handlers/...  # Unit tests only"
echo "   go test -v ./tests/integration/...     # Integration tests only"
echo "   go test -v ./... -run <TestName>         # Specific test"
