package main

import (
	"log"
	"proxy/internal/router"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	targetServers := router.NewTargetServers()
	proxyServer := router.NewProxyServer()
	var wg sync.WaitGroup

	for _, s := range targetServers {
		wg.Add(1)
		go func(s *router.Server) {
			defer wg.Done()
			if err := s.Run(); err != nil {
				log.Fatalf("Target Server exit: %v", err)
			}
		}(s)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := proxyServer.Run(); err != nil {
			log.Fatalf("Proxy Server exit: %v", err)
		}
	}()

	wg.Wait()
	log.Printf("Servers exited..")
}