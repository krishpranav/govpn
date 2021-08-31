package client

import (
	"log"
	"net"

	"github.com/krishpranav/govpn/common/cipher"
	"github.com/krishpranav/govpn/common/config"
	"github.com/krishpranav/govpn/vpn"
	"github.com/songgao/water/waterutil"
)

func StartUDPClient(config config.Config) {
	iface := vpn.CreateVpn(config.CIDR)
	serverAddr, err := net.ResolveUDPAddr("udp", config.ServerAddr)
	if err != nil {
		log.Fatalln("failed to resolve server addr:", err)
	}
	localAddr, err := net.ResolveUDPAddr("udp", config.LocalAddr)
	if err != nil {
		log.Fatalln("failed to get UDP socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed to listen on UDP socket:", err)
	}
	defer conn.Close()
	log.Printf("govpn udp client started on %v,CIDR is %v", config.LocalAddr, config.CIDR)

	go func() {
		buf := make([]byte, 1500)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil || n == 0 {
				continue
			}
			b := cipher.XOR(buf[:n])
			if !waterutil.IsIPv4(b) {
				continue
			}
			iface.Write(b)
		}
	}()

	packet := make([]byte, 1500)
	for {
		n, err := iface.Read(packet)
		if err != nil || n == 0 {
			continue
		}
		if !waterutil.IsIPv4(packet) {
			continue
		}
		b := cipher.XOR(packet[:n])
		conn.WriteToUDP(b, serverAddr)
	}
}
