package main

import (
	"log"
	"proxy/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	targetServer := router.NewTargetServer()
	proxyServer := router.NewProxyServer()

	go func() {
		if err := targetServer.Run(); err != nil {
			log.Fatalf("Error starting target server: %v", err)
		}
	}()

	if err := proxyServer.Run(); err != nil {
		log.Fatalf("Error starting proxy server: %v", err)
	}
}