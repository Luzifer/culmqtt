[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/culmqtt)](https://goreportcard.com/report/github.com/Luzifer/culmqtt)
![](https://badges.fyi/github/license/Luzifer/culmqtt)
![](https://badges.fyi/github/downloads/Luzifer/culmqtt)
![](https://badges.fyi/github/latest-release/Luzifer/culmqtt)

# Luzifer / culmqtt

`culmqtt` is a small connector between a CUL (in my case the [CC1101](http://busware.de/tiki-index.php?page=CUL)) and a MQTT server. It was written as I was migrating my home automation from FHEM to [Home-Assistent](https://www.home-assistant.io/) which lacks support for that CUL. Sadly some of my devices (for example the FS20 door-window-contacts) did not have an equivalent supported by Home-Assistent.

At first there was a shell command to forward Home-Assistent commands to FHEM which would execute them through the CUL. In order to shut down FHEM completely `culmqtt` replaced FHEM as the bridge between Home-Assistent and the CUL.
