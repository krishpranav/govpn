package config

import "github.com/krishpranav/govpn/common/cipher"

type Config struct {
	LocalAddr string
	ServerAddr string
	CIDR string
	Key string
	Protocol string
	ServerMode bool
}

func (config *Config) Init() {
	cipher.GenerateKey(config.Key)
}