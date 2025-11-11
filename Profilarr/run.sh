#!/usr/bin/env bash
set -e

CONFIG_PATH=/data/options.json

echo "Starting Profilarr..."

# Read configuration
export TZ=$(jq --raw-output '.TZ // "UTC"' $CONFIG_PATH)

echo "Configuration:"
echo "  TZ: $TZ"

# Set timezone
if [ -n "$TZ" ]; then
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime || true
    echo $TZ > /etc/timezone || true
fi

# Set config directory
export CONFIG_DIR=/config

# Ensure config directory exists
mkdir -p /config

# Change to app directory
cd /app

# Debug: List directory contents
echo "Contents of /app:"
ls -la /app/

# Check for possible entry points
if [ -f "server.js" ]; then
    echo "Starting Profilarr with server.js..."
    exec node server.js
elif [ -f "index.js" ]; then
    echo "Starting Profilarr with index.js..."
    exec node index.js
elif [ -f "backend/server.js" ]; then
    echo "Starting Profilarr backend..."
    exec node backend/server.js
elif [ -f "dist/server.js" ]; then
    echo "Starting Profilarr from dist..."
    exec node dist/server.js
else
    echo "ERROR: Could not find entry point!"
    echo "Available files:"
    find /app -type f -name "*.js" | head -20
    exit 1
fi
