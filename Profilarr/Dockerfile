ARG BUILD_FROM
FROM $BUILD_FROM

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install jq and su-exec (lightweight gosu alternative)
RUN apk add --no-cache jq su-exec

# Copy the run script
COPY run.sh /
RUN chmod a+x /run.sh

# Install Profilarr from official Docker image layers
COPY --from=santiagosayshey/profilarr:latest /app /app
COPY --from=santiagosayshey/profilarr:latest /entrypoint.sh /entrypoint.sh

# Create a symbolic link so gosu calls use su-exec
RUN ln -s /sbin/su-exec /usr/local/bin/gosu

WORKDIR /app

CMD [ "/run.sh" ]
