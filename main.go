package main

import (
	"flag"

	"github.com/krishpranav/govpn/client"
	"github.com/krishpranav/govpn/common/config"
	"github.com/krishpranav/govpn/server"
)

func main() {
	config := config.Config{}
	flag.StringVar(&config.CIDR, "c", "172.16.0.1/24", "tun interface CIDR")
	flag.StringVar(&config.LocalAddr, "l", "0.0.0.0:3000", "local address")
	flag.StringVar(&config.ServerAddr, "s", "0.0.0.0:3001", "server address")
	flag.StringVar(&config.Key, "k", "6w9z$C&F)J@NcRfWjXn3r4u7x!A%D*G-", "encryption key")
	flag.StringVar(&config.Protocol, "p", "wss", "protocol ws/wss/udp")
	flag.BoolVar(&config.ServerMode, "S", false, "server mode")
	flag.Parse()
	config.Init()
	switch config.Protocol {
	case "udp":
		if config.ServerMode {
			server.StartUDPServer(config)
		} else {
			client.StartUDPClient(config)
		}
	case "ws":
		if config.ServerMode {
			server.StartWSServer(config)
		} else {
			client.StartWSClient(config)
		}
	default:
	}
}
