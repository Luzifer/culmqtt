package main

import (
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func processFS20Message(housecode, device, command string) error {
	log.WithFields(log.Fields{
		"housecode": housecode,
		"device":    device,
		"command":   command,
	}).Info("FS20 status received")

	return publishFS20ToMQTT(housecode, device, command)
}

func publishFS20ToCUL(client mqtt.Client, msg mqtt.Message) {
	addr := strings.Split(msg.Topic(), "/")[1]
	cmd := string(msg.Payload())

	logger := log.WithFields(log.Fields{
		"address": addr,
		"command": cmd,
	})

	if _, err := fmt.Fprintf(port, "F%s%s\n", addr, cmd); err != nil {
		logger.WithError(err).Error("Unable to send message through CUL")
	}
	logger.Info("Message sent")
}

func publishFS20ToMQTT(housecode, device, command string) error {
	return errors.Wrap(
		mqttTokToErr(brokerClient.Publish(
			strings.Join([]string{"culmqtt", fmt.Sprintf("%s%s", housecode, device), "state"}, "/"),
			0x01, // QOS Level 1: At least once
			true,
			command,
		)),
		"publishing message",
	)
}
