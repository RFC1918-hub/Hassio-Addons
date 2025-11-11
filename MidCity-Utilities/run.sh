#!/usr/bin/env bash
set -e

echo "Starting MidCity Utilities Sensor..."

# Configuration is read by the Python script from /data/options.json
# The Python script handles validation

# Run the Python script
exec python3 /midcity_sensor.py
