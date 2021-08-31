package osutil

import (
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"

	"github.com/songgao/water"
)

func ConfigVpn(cidr string, iface *water.Interface) {
	os := runtime.GOOS
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Panicf("error cidr %v", cidr)
	}
	if os == "linux" {
		execCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "mtu", "1500")
		execCmd("/sbin/ip", "addr", "add", cidr, "dev", iface.Name())
		execCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "up")
	} else if os == "darwin" {
		minIp := ipNet.IP.To4()
		minIp[3]++
		execCmd("ifconfig", iface.Name(), "inet", ip.String(), minIp.String(), "up")
	} else if os == "windows" {
		log.Printf("please install openvpn client,see this link:%v", "https://github.com/OpenVPN/openvpn")
		log.Printf("open new cmd and enter:netsh interface ip set address name=\"%v\" source=static addr=%v mask=%v gateway=none", iface.Name(), ip.String(), ipNet.Mask.String())
	} else {
		log.Printf("not support os:%v", os)
	}
}
