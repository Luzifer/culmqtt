package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Luzifer/rconfig"
	"github.com/jacobsa/go-serial/serial"
	log "github.com/sirupsen/logrus"
)

var (
	cfg = struct {
		CULDevice      string `flag:"cul-device" default:"/dev/ttyACM0" env:"CUL_DEVICE" description:"TTY of the CUL to connect to"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		MQTTHost       string `flag:"mqtt-host" default:"tcp://127.0.0.1:1883" env:"MQTT_HOST" description:"Connection URI for the broker"`
		MQTTUser       string `flag:"mqtt-user" default:"" env:"MQTT_USER" description:"Username for broker connection"`
		MQTTPass       string `flag:"mqtt-pass" default:"" env:"MQTT_PASS" description:"Password for broker connection"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	port io.ReadWriteCloser

	version = "dev"
)

func init() {
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("culmqtt %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	options := serial.OpenOptions{
		PortName:        cfg.CULDevice,
		BaudRate:        19200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the port.
	var err error
	if port, err = serial.Open(options); err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	for {
		scanner := bufio.NewScanner(port)
		for scanner.Scan() {
			if err := processMessage(scanner.Text()); err != nil {
				log.WithError(err).Fatal("Unable to process message")
			}
		}
	}
}

func processMessage(message string) error {
	logger := log.WithField("message", message)

	switch message[0] {
	case 'F':
		return processFS20Message(
			message[1:5], // House code: 4 hex digits
			message[5:7], // Device code: 2 hex digits
			message[7:9], // Command: 2 hex digits
		)
	case 'V':
		// Version information, discard
		return nil
	default:
		logger.Error("Unknown message specifier")
		return nil
	}
}
