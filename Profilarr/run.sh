#!/usr/bin/env bash
set -e

CONFIG_PATH=/data/options.json

echo "Starting Profilarr..."

# Read configuration
export PUID=$(jq --raw-output '.PUID // 1000' $CONFIG_PATH)
export PGID=$(jq --raw-output '.PGID // 1000' $CONFIG_PATH)
export UMASK=$(jq --raw-output '.UMASK // "022"' $CONFIG_PATH)
export TZ=$(jq --raw-output '.TZ // "UTC"' $CONFIG_PATH)

echo "Configuration:"
echo "  PUID: $PUID"
echo "  PGID: $PGID"
echo "  UMASK: $UMASK"
echo "  TZ: $TZ"

# Start Profilarr with the original entrypoint
exec /init
