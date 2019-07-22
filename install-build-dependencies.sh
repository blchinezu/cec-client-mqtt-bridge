#!/bin/ash

apk add --no-cache --virtual=cec-client-mqtt-bridge-build-dependencies \
    git \
    go \
    musl-dev
