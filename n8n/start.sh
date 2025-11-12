#!/usr/bin/with-contenv bashio

# Read configuration from Home Assistant
N8N_HOST=$(bashio::config 'n8n_host')
N8N_PROTOCOL=$(bashio::config 'n8n_protocol')
TIMEZONE=$(bashio::config 'timezone')

# Set default values if not provided
if [ -z "$TIMEZONE" ]; then
    TIMEZONE="UTC"
fi

if [ -z "$N8N_PROTOCOL" ]; then
    N8N_PROTOCOL="http"
fi

# Export environment variables
export GENERIC_TIMEZONE="$TIMEZONE"
export TZ="$TIMEZONE"
export N8N_USER_FOLDER="/data"
export N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS="false"
export N8N_RUNNERS_ENABLED="true"
export N8N_TRUSTED_PROXY_IPS="*"
export NODE_ENV="production"

# Only set host/protocol/webhook if n8n_host is provided
if [ -n "$N8N_HOST" ]; then
    export N8N_HOST="$N8N_HOST"
    export N8N_PROTOCOL="$N8N_PROTOCOL"
    export WEBHOOK_URL="${N8N_PROTOCOL}://${N8N_HOST}"
    bashio::log.info "Configured n8n with host: ${N8N_HOST} (${N8N_PROTOCOL})"
else
    bashio::log.info "No n8n_host configured, using local access only"
fi

# Ensure data directory exists
mkdir -p /data

bashio::log.info "Starting n8n..."

# Start n8n
exec n8n start
