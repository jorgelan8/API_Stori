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


print_status "All tests completed successfully!"
echo ""
echo "🎉 Test Summary:"
echo "   ✅ Unit tests: PASSED"
echo "   ✅ Integration tests: PASSED"
echo ""
echo "📋 To run specific tests:"
echo "   go test -v ./internal/services/...     # Unit tests only"
echo "   go test -v ./tests/integration/...     # Integration tests only"
echo "   go test -v ./... -run <TestName>        # Specific test"
