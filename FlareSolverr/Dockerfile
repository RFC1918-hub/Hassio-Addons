ARG BUILD_FROM=ghcr.io/flaresolverr/flaresolverr:latest
FROM ${BUILD_FROM}

# Add Chrome flags to help with ARM64 and container environments
ENV CHROME_CLI="--disable-dev-shm-usage --no-sandbox --disable-setuid-sandbox --disable-gpu"
