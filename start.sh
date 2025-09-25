#!/bin/bash

# API Stori - Start Script
echo "🚀 Starting API Stori with Docker Compose..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

# Build and start the services
echo "📦 Building and starting services..."
#docker-compose up --build -d
docker-compose up

# Wait for the service to be ready
echo "⏳ Waiting for API to be ready..."
sleep 10

# Check if the API is responding
if curl -s http://localhost:8081/api/v1/health > /dev/null; then
    echo "✅ API is running successfully on port 8081!"
    echo ""
    echo "📊 Available endpoints:"
    echo "   POST http://localhost:8081/api/v1/migrate"
    echo "   GET  http://localhost:8081/api/v1/users/{user_id}/balance"
    echo "   GET  http://localhost:8081/api/v1/health"
    echo "   GET  http://localhost:8081/"
    echo ""
    echo "🧪 Test with sample CSV:"
    echo "   curl -X POST http://localhost:8081/api/v1/migrate -F \"csv_file=@examples/sample_transactions.csv\""
    echo ""
    echo "🔍 Check logs:"
    echo "   docker-compose logs -f"
    echo ""
    echo "🛑 Stop services:"
    echo "   docker-compose down"
else
    echo "❌ API is not responding. Check logs with: docker-compose logs"
    exit 1
fi
