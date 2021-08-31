#!bin/bash

sudo killall govpn-darwin-amd64
sudo ./bin/govpn-darwin-amd64 -S -l=:3001 -c=172.16.0.1/24 -p=wss &
echo "STARTED!!!!"