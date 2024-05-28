package config

import "log"

func Init(cfg *Config) {
    // TODO: Setup logger
    log.Printf("Logger initialized with config: %+v", cfg)
}