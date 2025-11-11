ARG BUILD_FROM=ghcr.io/fallenbagel/jellyseerr:latest
FROM ${BUILD_FROM}

USER root

# Set Jellyseerr to use /data for configuration
ENV CONFIG_DIRECTORY=/data

# Create startup script
RUN echo '#!/bin/sh' > /start.sh && \
    echo 'mkdir -p /data' >> /start.sh && \
    echo 'rm -rf /app/config' >> /start.sh && \
    echo 'ln -sf /data /app/config' >> /start.sh && \
    echo 'cd /app' >> /start.sh && \
    echo 'exec yarn start' >> /start.sh && \
    chmod +x /start.sh

CMD ["/start.sh"]
