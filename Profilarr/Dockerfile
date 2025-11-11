ARG BUILD_FROM=santiagosayshey/profilarr:latest
FROM ${BUILD_FROM}

# The base image already handles PUID/PGID, so we just need to ensure
# the config directory is properly set up

# Set Profilarr to use /data for configuration
ENV CONFIG_DIR=/data
