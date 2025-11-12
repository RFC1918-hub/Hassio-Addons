ARG BUILD_FROM=ghcr.io/hassio-addons/base:16.3.2
FROM docker.n8n.io/n8nio/n8n:latest as n8n

FROM ${BUILD_FROM}

# Copy n8n from the official image
COPY --from=n8n /usr/local/lib/node_modules/n8n /usr/local/lib/node_modules/n8n
COPY --from=n8n /usr/local/bin/n8n /usr/local/bin/n8n

# Install Node.js and npm
RUN apk add --no-cache \
    nodejs \
    npm \
    bash \
    jq

# Set n8n to use /data
ENV N8N_USER_FOLDER=/data
ENV N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS=false

# Copy start script
COPY start.sh /
RUN chmod a+x /start.sh

# Set the entrypoint
CMD ["/start.sh"]