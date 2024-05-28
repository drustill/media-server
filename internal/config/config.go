package config

import (
	"log"
	"net/http"
)

type Config struct {
	ServerAddress string
	Handler 		  func(w http.ResponseWriter, r *http.Request)
}

func load(addr string) *Config {
	if addr == "" {
			addr = ":8080"
	}
	log.Printf("Loaded configuration: addr=%s", addr)
	return &Config{
			ServerAddress: addr,
	}
}

func LoadTargetConfig(addr string) *Config {
	return load(addr)
}

func LoadProxyConfig(addr string) *Config {
	return load(addr)
}