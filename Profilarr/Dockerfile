ARG BUILD_FROM
FROM $BUILD_FROM

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install required packages
RUN apk add --no-cache \
    jq \
    nodejs \
    npm \
    python3 \
    py3-pip

# Copy the run script
COPY run.sh /
RUN chmod a+x /run.sh

# Install Profilarr from official Docker image layers
COPY --from=santiagosayshey/profilarr:latest /app /app

WORKDIR /app

CMD [ "/run.sh" ]
