ARG BUILD_FROM=docker.n8n.io/n8nio/n8n:latest
FROM ${BUILD_FROM}

USER root

# Set n8n to use /data
ENV N8N_USER_FOLDER=/data
ENV N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS=false

# Create entrypoint that ensures /data exists
ENTRYPOINT ["/bin/sh", "-c", "mkdir -p /data && n8n start"]