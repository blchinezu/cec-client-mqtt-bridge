package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "regexp"

    MQTT "github.com/eclipse/paho.mqtt.golang"
)

func mqtt2cecHandler(client MQTT.Client, msg MQTT.Message) {
    fmt.Printf("%s\n", msg.Payload())
}

func main() {

    // Config MQTT
    opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("cec-client-mqtt")
    opts.SetCleanSession(true)

    // Connect
    mqttClient := MQTT.NewClient(opts)
    if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    // MQTT to CEC-CLIENT
    if token := mqttClient.Subscribe("cec/client/tx", 0, mqtt2cecHandler); token.Wait() && token.Error() != nil {
        fmt.Println(token.Error())
        os.Exit(1)
    }

    // CEC-CLIENT to MQTT
    stdin := bufio.NewReader(os.Stdin)

    // Regex to get message topic & contents
    rTopicContents := regexp.MustCompile(`^([A-Z]*)\:\s*\[.{18}(.*)`)

    // Init message to be published
    for {

        // Get CEC message
        cecMessage, err := stdin.ReadString('\n')

        // Check for errors
        if err == io.EOF {
            os.Exit(0)
        }

        // Find topic & contents
        matches := rTopicContents.FindStringSubmatch(cecMessage)

        switch len(matches) {
        case 0:
            mqttClient.Publish("cec/client/rx/UNKNOWN", 0, false, cecMessage)
        case 3:
            if matches[1] == "TRAFFIC" {
                mqttClient.Publish("cec/client/rx/"+matches[1], 0, false, matches[2])
            }
        }
    }
}
