package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/krishpranav/govpn/common/cipher"
	"github.com/krishpranav/govpn/common/config"
	"github.com/krishpranav/govpn/common/netutil"
	"github.com/krishpranav/govpn/vpn"
	"github.com/patrickmn/go-cache"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

func StartUDPServer(config config.Config) {
	iface := vpn.CreateVpn(config.CIDR)
	localAddr, err := net.ResolveUDPAddr("udp", config.LocalAddr)
	if err != nil {
		log.Fatalln("failed to get UDP socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed to listen on UDP socket:", err)
	}
	defer conn.Close()
	log.Printf("govpn udp server started on %v,CIDR is %v", config.LocalAddr, config.CIDR)

	forwarder := &Forwarder{localConn: conn, connCache: cache.New(30*time.Minute, 10*time.Minute)}
	go forwarder.forward(iface, conn)

	buf := make([]byte, 1500)
	for {
		n, cliAddr, err := conn.ReadFromUDP(buf)
		if err != nil || n == 0 {
			continue
		}
		b := cipher.XOR(buf[:n])
		if !waterutil.IsIPv4(b) {
			continue
		}
		iface.Write(b)
		srcAddr, dstAddr := netutil.GetAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		key := fmt.Sprintf("%v->%v", srcAddr, dstAddr)
		forwarder.connCache.Set(key, cliAddr, cache.DefaultExpiration)
	}
}

type Forwarder struct {
	localConn *net.UDPConn
	connCache *cache.Cache
}

func (f *Forwarder) forward(iface *water.Interface, conn *net.UDPConn) {
	packet := make([]byte, 1500)
	for {
		n, err := iface.Read(packet)
		if err != nil || n == 0 {
			continue
		}
		b := packet[:n]
		if !waterutil.IsIPv4(b) {
			continue
		}
		srcAddr, dstAddr := netutil.GetAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		key := fmt.Sprintf("%v->%v", dstAddr, srcAddr)
		v, ok := f.connCache.Get(key)
		if ok {
			b = cipher.XOR(b)
			f.localConn.WriteToUDP(b, v.(*net.UDPAddr))
		}
	}
}
