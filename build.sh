#!/bin/ash

go get github.com/eclipse/paho.mqtt.golang
go build -ldflags "-s -w" src/cec-client-mqtt-bridge.go
