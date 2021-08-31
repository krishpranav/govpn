# govpn
A Simple VPN Built In Golang

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

# Installation
```
$ git clone https://github.com/krishpranav/govpn
```

# Usage:
```
Usage of ./main:
  -S    server mode
  -c string
        vpn interface CIDR (default "172.16.0.1/24")
  -k string
        encryption key (default "S#Q#FBSDAE#%!@#!@#%!NDADSA")
  -p string
        protocol ws/wss/udp (default "wss")
  -l string
        local address (default "0.0.0.0:3000")
  -s string
        server address (default "0.0.0.0:3001")        
```

# Build:
```
$ bash scripts/build.sh
```

# Server:
```
sudo ./main -S -l=:3001 -c=172.16.0.1/24 -k=123456
```

# Client:
```
sudo ./main -l=:3000 -s=server-addr:3001 -c=172.16.0.10/24 -k=123456
```

# Server Setup:

- Add TLS for websocket,reverse proxy server(3001) via nginx/caddy(443)

- Enable IP forwarding on server

```
  sudo echo 1 > /proc/sys/net/ipv4/ip_forward
  sudo sysctl -p
  sudo iptables -t nat -A POSTROUTING -s 172.16.0.0/24 -o ens3 -j MASQUERADE
  sudo apt-get install iptables-persistent
  sudo iptables-save > /etc/iptables/rules.v4
```
