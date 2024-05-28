package router

import (
	"log"
	"net/http"
	"os"
	"proxy/internal/config"
	"strings"
	"time"
)

type Server struct {
	name			 string
	httpServer *http.Server
}

func NewTargetServers() []*Server {
	var servers []*Server
	targetAddrs := strings.Split(os.Getenv("TARGET_ADDRS"), ",")
	for _, addr := range targetAddrs {
		servers = append(servers, new(config.LoadTargetConfig(addr), ServerHandler, addr[8:]))
	}
	return servers
}

func NewProxyServer() *Server {
	addr := os.Getenv("PROXY_ADDRESS")
	cfg := config.LoadProxyConfig(addr)
	lbKind := os.Getenv("LOADBALANCER_KIND")
	lb := NewLoadBalancer(lbKind, cfg.TargetAddresses)
	p := NewProxy(lb)
	return new(cfg, p.ServeHTTP, "Proxy")
}

func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] [Hit: %s %s]", s.name, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("[%s] [Complete %s] [RTT: %v]", s.name, r.URL.Path, time.Since(start))
})
}

func new(cfg *config.Config, f http.HandlerFunc, name string) *Server {
	mux := http.NewServeMux()

	s := &Server{
			name: name,
			httpServer: &http.Server{
					Addr:    cfg.ServerAddress,
					Handler: mux,
			},
	}
	mux.Handle("/", s.middleware(f))

	return s
}

func (s *Server) Run() error {
	log.Printf("Starting %s Server on %s", s.name, s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}