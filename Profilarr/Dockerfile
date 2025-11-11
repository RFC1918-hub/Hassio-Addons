ARG BUILD_FROM
FROM $BUILD_FROM

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install dependencies
RUN apk add --no-cache \
    bash \
    curl \
    nodejs \
    npm \
    git

# Install Profilarr
RUN git clone https://github.com/santiagosayshey/profilarr.git /app && \
    cd /app && \
    npm install

# Copy the run script
COPY run.sh /
RUN chmod a+x /run.sh

WORKDIR /app

CMD [ "/run.sh" ]
