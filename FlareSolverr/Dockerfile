ARG BUILD_FROM=ghcr.io/flaresolverr/flaresolverr:latest
FROM ${BUILD_FROM}

USER root

# Patch utils.py to add Chrome flags for ARM64/container compatibility
RUN sed -i '/options\.add_argument.*headless/a \    options.add_argument("--disable-dev-shm-usage")\n    options.add_argument("--no-sandbox")\n    options.add_argument("--disable-setuid-sandbox")\n    options.add_argument("--disable-gpu")\n    options.add_argument("--disable-software-rasterizer")' /app/utils.py
