package main

import (
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var brokerClient mqtt.Client

func init() {
	opts := mqtt.NewClientOptions().AddBroker(cfg.MQTTHost)
	if cfg.MQTTUser != "" || cfg.MQTTPass != "" {
		opts.SetUsername(cfg.MQTTUser).SetPassword(cfg.MQTTPass)
	}

	brokerClient = mqtt.NewClient(opts)

	if err := mqttTokToErr(brokerClient.Connect()); err != nil {
		log.WithError(err).Fatal("Connect to MQTT broker")
	}

	if err := mqttTokToErr(brokerClient.Subscribe("culmqtt/+/send", 0x01, publishFS20ToCUL)); err != nil {
		log.WithError(err).Fatal("Subscribe to topic")
	}
}

func mqttTokToErr(tok mqtt.Token) error {
	if !tok.WaitTimeout(cfg.MQTTTimeout) {
		return errors.New("command timed out")
	}

	return tok.Error()
}
