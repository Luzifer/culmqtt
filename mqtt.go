package main

import mqtt "github.com/eclipse/paho.mqtt.golang"

var brokerClient mqtt.Client

func init() {
	opts := mqtt.NewClientOptions().AddBroker(cfg.MQTTHost)
	if cfg.MQTTUser != "" || cfg.MQTTPass != "" {
		opts.SetUsername(cfg.MQTTUser).SetPassword(cfg.MQTTPass)
	}

	brokerClient = mqtt.NewClient(opts)

	brokerClient.Connect().Wait()
	brokerClient.Subscribe("culmqtt/+/send", 0x01, publishFS20ToCUL)
}
