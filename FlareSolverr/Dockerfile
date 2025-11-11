ARG BUILD_FROM=ghcr.io/flaresolverr/flaresolverr:latest
FROM ${BUILD_FROM}

# Create a wrapper script for Chromium with necessary flags for ARM64/containers
RUN mv /usr/bin/chromium /usr/bin/chromium-browser && \
    echo '#!/bin/bash' > /usr/bin/chromium && \
    echo 'exec /usr/bin/chromium-browser \
      --disable-dev-shm-usage \
      --no-sandbox \
      --disable-setuid-sandbox \
      --disable-gpu \
      --disable-software-rasterizer \
      --disable-extensions \
      "$@"' >> /usr/bin/chromium && \
    chmod +x /usr/bin/chromium
