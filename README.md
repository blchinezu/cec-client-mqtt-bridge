# cec-client-mqtt-bridge
HDMI CEC-MQTT Bridge using the **cec-client** binary

--------------------------------------------------------------------------------

### Description

This lets you send commands and receive answers from the `cec-client` binary
through MQTT.

 - To send commands to `cec-client` you publish into `cec/client/tx`
 - To listen for `cec-client` output you subscribe to `cec/client/rx/TRAFFIC`

--------------------------------------------------------------------------------

### Why this exists

All the existing implementations try to connect directly to CEC and interpret
the messages although this is already done in `cec-client` and it's the most
reliable implementation. Everyone checks the CEC functionality and/or their
implementation with `cec-client` anyway.

So this project is creating a bridge between the `cec-client` binary and an MQTT
broker placed on localhost. Basically it's creating a dual-way pipe between the
services.

--------------------------------------------------------------------------------

### How's this working

 - The `cec-client-mqtt-bridge` binary connects to the MQTT broker
 - It subscribes to `cec/client/tx` on a different thread
 - It spills the received messages in `stdout` which is piped to `cec-client`
 - It then listens for `stdin` input which is piped from `cec-client` on the
   main thread
 - When receiving something it publishes the message on `cec/client/rx/TRAFFIC`

--------------------------------------------------------------------------------

### Hardware

Make sure you have HDMI CEC capable hardware

 - Raspberry PIs have this by default
 - Intel NUCs usually require an [additional hardware module](https://www.pulse-eight.com/p/154/intel-nuc-hdmi-cec-adapter) to be installed
 - Other devices can use [this generic adapter](https://www.pulse-eight.com/p/104/usb-hdmi-cec-adapter)

--------------------------------------------------------------------------------

### Build

The build/dependency/run scripts are built for Apline Linux but can be easily
adapted for Debian or something else.

```bash
# Install git, go & musl-dev
ash install-build-dependencies.sh

# Build binary
ash build.sh

# Remove git, go & musl-dev
ash remove-build-dependencies.sh
```

--------------------------------------------------------------------------------

### Runtime dependencies

For this to work you need the `cec-client` binary.

To install it in Alpine Linux just run:

```bash
apk add libcec
```

For Debian based systems this should do the trick:

```bash
apt install cec-utils
```

--------------------------------------------------------------------------------

### Run

Creates two pipes in `/tmp/fifo/` which will facilitate the communication
between processes and then launches the processes.

```bash
bash run.sh
```

Here's how `run.sh` looks like:

```bash
# Create pipes
mkdir -p /tmp/fifo
rm -rf /tmp/fifo/mqtt2cec /tmp/fifo/cec2mqtt
mkfifo /tmp/fifo/mqtt2cec /tmp/fifo/cec2mqtt

# Launch bridge
./cec-client-mqtt-bridge > /tmp/fifo/mqtt2cec < /tmp/fifo/cec2mqtt &

# Launch cec client
cec-client < /tmp/fifo/mqtt2cec > /tmp/fifo/cec2mqtt

```

--------------------------------------------------------------------------------

### Docker

To use this in an Alpine based Docker image you could use this:

```dockerfile
ENV CECCLIENT_MQTT_BRIDGE_PATH "/app/modules/cec-client-mqtt-bridge"
RUN echo -e "\n > INSTALL CEC-CLIENT-MQTT-BRIDGE\n" \
 && apk add --no-cache --virtual=build-dependencies \
    git \
    \
 && mkdir -p $CECCLIENT_MQTT_BRIDGE_PATH/src \
 && cd $CECCLIENT_MQTT_BRIDGE_PATH \
 && git clone https://github.com/blchinezu/cec-client-mqtt-bridge.git ./src \
 && ash src/install-build-dependencies.sh \
 && ash src/build.sh \
 \
 && echo -e "\n > CLEANUP\n" \
 && ash src/remove-build-dependencies.sh \
 && apk del --purge build-dependencies \
 && rm -rf \
    ./src \
    /root/.cache \
    /root/go \
    /tmp/* \
    /var/tmp/*
```
