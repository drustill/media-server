package router

import (
	"testing"
)

func TestRoundRobinBalancer(t *testing.T) {
	t.Run("it should circularly select servers", func(t *testing.T) {
		inputServers := []string{"localhost:9090", "localhost:9091"}
		servers := []string{"http://localhost:9090", "http://localhost:9091"}
		lb := NewRoundRobinBalancer(inputServers)

		if lb.SelectTarget() != servers[0] {
				t.Errorf("Expected %s, got %s", servers[0], lb.SelectTarget())
		}
		if lb.SelectTarget() != servers[1] {
				t.Errorf("Expected %s, got %s", servers[1], lb.SelectTarget())
		}
		if lb.SelectTarget() != servers[0] {
				t.Errorf("Expected %s, got %s", servers[0], lb.SelectTarget())
		}
	})
}

func TestLeastActiveBalancer(t *testing.T) {
	t.Run("it should select the least active server", func(t *testing.T) {
		inputServers := []string{"localhost:9090", "localhost:9091"}
		lb := NewLeastConnectionBalancer(inputServers)

		server := lb.SelectTarget()
		lb.RecordRequest()

		nextServer := lb.SelectTarget()
		if server == nextServer {
				t.Errorf("Expected different server, got same server")
		}
	})

	t.Run("it should keep track of active connections", func(t *testing.T) {
		inputServers := []string{"localhost:9090", "localhost:9091"}
		lb := NewLeastConnectionBalancer(inputServers)
		
		lb.SelectTarget() // Select first server
		lb.RecordRequest() // Increment active connections
		lb.SelectTarget() // Shift to next server
		lb.RecordRequest() // Increment active connections

		// Both servers should have one active connection
		for _, s := range lb.active {
			if s != 1 {
				t.Errorf("Expected active connections, got 0")
			}
		}
	})

	t.Run("it should decrement active connections", func(t *testing.T) {
		inputServers := []string{"localhost:9090", "localhost:9091"}
		lb := NewLeastConnectionBalancer(inputServers)
		
		s := lb.SelectTarget() // Select first server
		lb.RecordRequest() // Increment active connections
		s = lb.SelectTarget() // Shift to next server
		lb.RecordRequest() // Increment active connections
		lb.RecordResponse(s) // Decrement active connections

		// At least one server should have no active connections
		if lb.active[0] != 0 && lb.active[1] != 0 {
			t.Errorf("Expected active connections, got 1")
		}
	})

	t.Run("it should decrement the server that responded", func(t *testing.T) {
		inputServers := []string{"localhost:9090", "localhost:9091"}
		lb := NewLeastConnectionBalancer(inputServers)
		
		s := lb.SelectTarget() // Select first server
		lb.RecordRequest() // Increment active connections
		lb.SelectTarget() // Shift to next server
		lb.RecordRequest() // Increment active connections

		lb.RecordResponse(s) // Should decrement the first target

		if lb.active[0] != 0 {
			t.Errorf("Expected active connections, got 1")
		}
	})
}

