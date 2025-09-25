#!/bin/bash

# API Stori - Stop Script
echo "🛑 Stopping API Stori..."

# Stop and remove containers
docker-compose down

echo "✅ API Stori stopped successfully!"
