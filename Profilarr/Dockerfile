ARG BUILD_FROM=santiagosayshey/profilarr:latest
FROM $BUILD_FROM

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Copy the run script
COPY run.sh /
RUN chmod a+x /run.sh

WORKDIR /

CMD [ "/run.sh" ]
