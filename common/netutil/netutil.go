package netutil

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/gorilla/websocket"
	"github.com/net-byte/vtun/common/config"
	"github.com/songgao/water/waterutil"
)

func ConnectWS(config config.Config) *websocket.Conn {
	scheme := "ws"
	if config.Protocol == "wss" {
		scheme = "wss"
	}
	u := url.URL{Scheme: scheme, Host: config.ServerAddr, Path: "/way-to-freedom"}
	header := make(http.Header)
	header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Printf("[client] failed to dial websocket %v", err)
		return nil
	}
	return c
}

func CloseWS(wsConn *websocket.Conn) {
	wsConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second*5))
	wsConn.Close()
}
