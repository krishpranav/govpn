package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/krishpranav/govpn/common/cipher"
	"github.com/krishpranav/govpn/common/config"
	"github.com/krishpranav/govpn/common/netutil"
	"github.com/krishpranav/govpn/register"
	"github.com/krishpranav/govpn/vpn"
	"github.com/patrickmn/go-cache"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1500,
	WriteBufferSize:   1500,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StartWSServer start ws server
func StartWSServer(config config.Config) {
	iface := vpn.CreateVpn(config.CIDR)
	c := cache.New(30*time.Minute, 10*time.Minute)
	go vpnToWs(iface, c)
	log.Printf("govpn ws server started on %v,CIDR is %v", config.LocalAddr, config.CIDR)
	http.HandleFunc("/way-to-freedom", func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		wsToVpn(wsConn, iface, c)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello，世界！")
	})

	http.HandleFunc("/ip", func(w http.ResponseWriter, req *http.Request) {
		ip := req.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = strings.Split(req.RemoteAddr, ":")[0]
		}
		resp := fmt.Sprintf("%v", ip)
		io.WriteString(w, resp)
	})

	http.HandleFunc("/register/pick/ip", func(w http.ResponseWriter, req *http.Request) {
		key := req.Header.Get("key")
		if key != config.Key {
			error403(w, req)
			return
		}
		ip, pl := register.PickClientIP(config.CIDR)
		resp := fmt.Sprintf("%v/%v", ip, pl)
		io.WriteString(w, resp)
	})

	http.HandleFunc("/register/delete/ip", func(w http.ResponseWriter, req *http.Request) {
		key := req.Header.Get("key")
		if key != config.Key {
			error403(w, req)
			return
		}
		ip := req.URL.Query().Get("ip")
		if ip != "" {
			register.DeleteClientIP(ip)
		}
		io.WriteString(w, "OK")
	})

	http.HandleFunc("/register/keepalive/ip", func(w http.ResponseWriter, req *http.Request) {
		key := req.Header.Get("key")
		if key != config.Key {
			error403(w, req)
			return
		}
		ip := req.URL.Query().Get("ip")
		if ip != "" {
			register.KeepAliveClientIP(ip)
		}
		io.WriteString(w, "OK")
	})

	http.HandleFunc("/register/list/ip", func(w http.ResponseWriter, req *http.Request) {
		key := req.Header.Get("key")
		if key != config.Key {
			error403(w, req)
			return
		}
		io.WriteString(w, strings.Join(register.ListClientIP(), "\r\n"))
	})

	http.ListenAndServe(config.LocalAddr, nil)
}

func error403(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("403 No Permission"))
}

func vpnToWs(iface *water.Interface, c *cache.Cache) {
	buffer := make([]byte, 1500)
	for {
		n, err := iface.Read(buffer)
		if err != nil || err == io.EOF || n == 0 {
			continue
		}
		b := buffer[:n]
		if !waterutil.IsIPv4(b) {
			continue
		}
		srcAddr, dstAddr := netutil.GetAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		key := fmt.Sprintf("%v->%v", dstAddr, srcAddr)
		v, ok := c.Get(key)
		if ok {
			b = cipher.XOR(b)
			v.(*websocket.Conn).WriteMessage(websocket.BinaryMessage, b)
		}
	}
}

func wsToVpn(wsConn *websocket.Conn, iface *water.Interface, c *cache.Cache) {
	defer netutil.CloseWS(wsConn)
	for {
		wsConn.SetReadDeadline(time.Now().Add(time.Duration(30) * time.Second))
		_, b, err := wsConn.ReadMessage()
		if err != nil || err == io.EOF {
			break
		}
		b = cipher.XOR(b)
		if !waterutil.IsIPv4(b) {
			continue
		}
		srcAddr, dstAddr := netutil.GetAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		key := fmt.Sprintf("%v->%v", srcAddr, dstAddr)
		c.Set(key, wsConn, cache.DefaultExpiration)
		iface.Write(b[:])
	}
}
