#!/bin/bash

# Create pipes
mkdir -p /tmp/fifo
rm -rf /tmp/fifo/mqtt2cec /tmp/fifo/cec2mqtt
mkfifo /tmp/fifo/mqtt2cec /tmp/fifo/cec2mqtt

# Launch bridge
./cec-client-mqtt-bridge > /tmp/fifo/mqtt2cec < /tmp/fifo/cec2mqtt &

# Launch cec client
cec-client < /tmp/fifo/mqtt2cec > /tmp/fifo/cec2mqtt
