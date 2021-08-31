#!bin/bash
export GO111MODULE=on

# Linux Build
GOOS=linux GOARCH=amd64 go build -o ./bin/govpn-linux-amd64 ./main.go

# LinuxARM Build
GOOS=linux GOARCH=arm64 go build -o ./bin/govpn-linux-arm64 ./main.go

# macOS Build
GOOS=darwin GOARCH=amd64 go build -o ./bin/govpn-darwin-amd64 ./main.go

# Windows Build
GOOS=windows GOARCH=amd64 go build -o ./bin/govpn-windows-amd64.exe ./main.go

echo "DONE!!!"