ARG BUILD_FROM
FROM $BUILD_FROM

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install jq for parsing config
RUN apk add --no-cache jq

# Copy the run script
COPY run.sh /
RUN chmod a+x /run.sh

# Install Profilarr from official Docker image layers
COPY --from=santiagosayshey/profilarr:latest /app /app
COPY --from=santiagosayshey/profilarr:latest /entrypoint.sh /entrypoint.sh

WORKDIR /app

CMD [ "/run.sh" ]
