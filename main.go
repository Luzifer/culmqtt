package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Luzifer/rconfig/v2"
	"github.com/jacobsa/go-serial/serial"
	log "github.com/sirupsen/logrus"
)

var (
	cfg = struct {
		CULDevice      string        `flag:"cul-device" default:"/dev/ttyACM0" description:"TTY of the CUL to connect to"`
		LogLevel       string        `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		MQTTHost       string        `flag:"mqtt-host" default:"tcp://127.0.0.1:1883" description:"Connection URI for the broker"`
		MQTTUser       string        `flag:"mqtt-user" default:"" description:"Username for broker connection"`
		MQTTPass       string        `flag:"mqtt-pass" default:"" description:"Password for broker connection"`
		MQTTTimeout    time.Duration `flag:"mqtt-timeout" default:"2s" description:"Timeout for MQTT actions"`
		VersionAndExit bool          `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	port io.ReadWriteCloser

	version = "dev"
)

func init() {
	rconfig.AutoEnv(true)
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

	// Send initialization for the CUL
	// TODO: This might be useful to be configurable?
	fmt.Fprintln(port, "Ax")  // reset AskSin
	fmt.Fprintln(port, "Zx")  // reset Moritz
	fmt.Fprintln(port, "brx") // reset WMBus
	fmt.Fprintln(port, "X21") // Turn on echoing of received messages

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

	if message == "" {
		return nil
	}

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
