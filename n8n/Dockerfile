ARG BUILD_FROM=ghcr.io/hassio-addons/base:16.3.2
FROM docker.n8n.io/n8nio/n8n:latest AS n8n

FROM ${BUILD_FROM}

# Switch to root to install packages
USER root

# Copy n8n from the official image
COPY --from=n8n /usr/local/lib/node_modules/n8n /usr/local/lib/node_modules/n8n
COPY --from=n8n /usr/local/bin/n8n /usr/local/bin/n8n

# Install Node.js and npm
RUN apk add --no-cache \
    nodejs \
    npm

# Set n8n to use /data
ENV N8N_USER_FOLDER=/data
ENV N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS=false

# Create S6 service directory and copy run script
RUN mkdir -p /etc/services.d/n8n
COPY run.sh /etc/services.d/n8n/run
RUN chmod a+x /etc/services.d/n8n/run