#!/usr/bin/with-contenv bashio

# Print nice header
bashio::log.info "Starting MidCity Utilities Sensor..."

# Read configuration from options
USERNAME=$(bashio::config 'username')
PASSWORD=$(bashio::config 'password')
SCAN_INTERVAL=$(bashio::config 'scan_interval')

# Validate configuration
if bashio::var.is_empty "${USERNAME}"; then
    bashio::exit.nok "Username is required but not provided in configuration"
fi

if bashio::var.is_empty "${PASSWORD}"; then
    bashio::exit.nok "Password is required but not provided in configuration"
fi

bashio::log.info "Configuration validated successfully"
bashio::log.info "Scan interval: ${SCAN_INTERVAL:-300} seconds"

# Run the Python script
python3 /midcity_sensor.py
