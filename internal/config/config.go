package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	ServerAddress string
	TargetAddresses []string
}

func load(src string, dsts []string) *Config {
	if src == "" {
			src = ":8080"
	}
	return &Config{
			ServerAddress: src,
			TargetAddresses: dsts,
	}
}

func LoadTargetConfig(addr string) *Config {
	return load(addr, []string{})
}

func LoadProxyConfig(addr string) *Config {
	targetDsts := strings.Split(os.Getenv("TARGET_ADDRS"), ",")
	if len(targetDsts) == 0 {
			log.Fatal("TARGET_ADDRS is not set or is empty")
	}

	return load(addr, targetDsts)
}