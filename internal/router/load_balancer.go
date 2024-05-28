package router

import (
	"slices"
	"sync/atomic"
)


type LoadBalancer interface {
	SelectTarget() 	 						string   		// Transition to the next target
	RecordRequest()  									 			// Least Connection Algorithm
	RecordResponse(s string) 								// Least Connection Algorithm
}

func NewLoadBalancer(lb string, servers []string) LoadBalancer {
	switch lb {
	case "round_robin":
		return NewRoundRobinBalancer(servers)
	case "least_connection":
		return NewLeastConnectionBalancer(servers)
	default:
		return NewRoundRobinBalancer(servers)
	}
}

type LeastConnectionBalancer struct {
	servers []string
	curr    uint32
	active	[]uint32
}

func NewLeastConnectionBalancer(servers []string) *LeastConnectionBalancer {
	baseUrls := make([]string, len(servers))
	for i, s := range servers {
			baseUrls[i] = "http://" + s
	}
	return &LeastConnectionBalancer{
			servers: baseUrls,
			curr:    0,
			active: make([]uint32, len(servers)), // Go by index
	}
}

func (l *LeastConnectionBalancer) SelectTarget() string {
	l.curr = uint32(slices.Index(l.active, slices.Min(l.active)))	
	return l.servers[l.curr]
}

func (l *LeastConnectionBalancer) RecordRequest() {
	atomic.AddUint32(&l.active[l.curr], 1)
}
func (l *LeastConnectionBalancer) RecordResponse(s string) {
	sIndex := slices.Index(l.servers, s)
	atomic.AddUint32(&l.active[sIndex], ^uint32(0)) // Decrement
}


type RoundRobinBalancer struct {
	servers []string
	curr   uint32
}

func NewRoundRobinBalancer(servers []string) *RoundRobinBalancer {
	baseUrls := make([]string, len(servers))
	for i, s := range servers {
			baseUrls[i] = "http://" + s
	}
	return &RoundRobinBalancer{
			servers: baseUrls,
			curr:   0,
	}
}

func (r *RoundRobinBalancer) SelectTarget() string {
	i := atomic.AddUint32(&r.curr, 1) // Atomic is thread-safe
	return r.servers[(i-1) % uint32(len(r.servers))]
}
// No-op for RoundRobinBalancer
func (r *RoundRobinBalancer) RecordRequest() {}
func (r *RoundRobinBalancer) RecordResponse(s string) {}