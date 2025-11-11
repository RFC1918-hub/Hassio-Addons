ARG BUILD_FROM=santiagosayshey/profilarr:latest
FROM ${BUILD_FROM}

USER root

# Set Profilarr to use /data for configuration
ENV CONFIG_DIR=/config

# Create startup script
RUN echo '#!/bin/sh' > /start.sh && \
    echo 'mkdir -p /data' >> /start.sh && \
    echo 'rm -rf /config' >> /start.sh && \
    echo 'ln -sf /data /config' >> /start.sh && \
    echo 'exec /app/profilarr' >> /start.sh && \
    chmod +x /start.sh

CMD ["/start.sh"]
