ARG BUILD_FROM
FROM $BUILD_FROM

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install requirements for add-on
RUN \
  apk add --no-cache \
    python3 \
    py3-pip

# Copy requirements
COPY requirements.txt /

# Install Python dependencies
RUN pip3 install --no-cache-dir -r /requirements.txt

# Copy data for add-on
COPY run.sh /
COPY midcity_sensor.py /

RUN chmod a+x /run.sh

CMD [ "/run.sh" ]
