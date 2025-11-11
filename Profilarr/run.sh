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

# Start the backend
echo "Starting Profilarr backend and frontend..."
exec node backend/server.js
