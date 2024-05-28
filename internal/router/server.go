package router

import (
	"log"
	"net/http"
	"os"
	"proxy/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewTargetServer() *Server {
	addr := os.Getenv("TARGET_ADDRESS")
	cfg := config.LoadTargetConfig(addr)
	return new(cfg, ServerHandler)
}

func NewProxyServer() *Server {
	addr := os.Getenv("PROXY_ADDRESS")
	cfg := config.LoadProxyConfig(addr)
	return new(cfg, ProxyHandler)
}

func new(cfg *config.Config, f func(w http.ResponseWriter, r *http.Request)) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", f)

	return &Server{
			httpServer: &http.Server{
					Addr:    cfg.ServerAddress,
					Handler: mux,
			},
	}
}

func (s *Server) Run() error {
	log.Printf("Starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}