ARG BUILD_FROM=ghcr.io/flaresolverr/flaresolverr:latest
FROM ${BUILD_FROM}

USER root

# Create a wrapper script for Chromium with necessary flags for ARM64/containers
RUN echo '#!/bin/sh' > /usr/local/bin/chromium-wrapper && \
    echo 'exec /usr/bin/chromium \
      --disable-dev-shm-usage \
      --no-sandbox \
      --disable-setuid-sandbox \
      --disable-gpu \
      --disable-software-rasterizer \
      --disable-extensions \
      "$@"' >> /usr/local/bin/chromium-wrapper && \
    chmod +x /usr/local/bin/chromium-wrapper

# Set environment variable to use the wrapper
ENV CHROME_EXE_PATH=/usr/local/bin/chromium-wrapper
