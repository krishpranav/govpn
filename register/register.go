package register

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var _register *cache.Cache

func init() {
	_register = cache.New()
}

func AddClientIP(ip string) {
	_register.Add(ip, 0, cache.DefaultExpiration)
}

func ExistClientIP(ip string) bool {
	_, ok := _register.Get(ip)
	return ok
}

func KeepAliveClientIP(ip string) {
	if ExistsClientIP(ip) {
		_register.Increment(ip, 1)
	} else {
		AddClientIP(ip)
	}
}

func checkIPv4(ip net.IP) net.IP {
	if v4 := ip.To4(); v4 != nil {
		return v4
	}

	return ip
}
