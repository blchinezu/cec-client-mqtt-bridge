# cec-client-mqtt-bridge
HDMI CEC-MQTT Bridge using the **cec-client** binary

--------------------------------------------------------------------------------

### Description

All the existing implementations try to connect directly to CEC and interpret the messages although this is already done in cec-client and it's the most reliable implementation. Everyone checks the CEC functionality with cec-client anyway.

So this project is creating a bridge between the cec-client binary and an MQTT broker placed on localhost. Basically it's creating a dual-way pipe between the services.

The output of MQTT's topic **cec/client/tx** is going into cec-client and the output of cec-client is going into **cec/client/rx/TRAFFIC**

--------------------------------------------------------------------------------

### Scripts

The build/dependency/run scripts are built for Apline Linux. But can be easily adapted for other systems.

--------------------------------------------------------------------------------

### Install build dependencies

```bash
ash install-build-dependencies.sh
```
--------------------------------------------------------------------------------

### Build

```bash
ash build.sh
```
--------------------------------------------------------------------------------

### Remove build dependencies

```bash
ash remove-build-dependencies.sh
```
--------------------------------------------------------------------------------

### Run

```bash
bash run.sh
```

This creates two pipes in `/tmp/fifo/` which will facilitate the communication between processes. It's then launching the processes.

Here's how the `run.sh` script looks like:

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