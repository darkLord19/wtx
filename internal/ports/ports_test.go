package ports

import (
	"fmt"
	"net"
	"testing"
)

func getFreePort(t *testing.T) int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to resolve tcp addr: %v", err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to listen on port 0: %v", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func TestIsPortOpen(t *testing.T) {
	// 1. Get a free port but don't hold it (so it's closed/free)
	port := getFreePort(t)

	// Currently free, so IsPortOpen (In Use) should be false
	if IsPortOpen(port) {
		t.Errorf("Port %d should be free (IsPortOpen=false), but got true", port)
	}

	// 2. Occupy a port
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// If we can't listen, maybe it was taken in the microsecond between. Retry?
		// For unit test, we'll fail.
		t.Fatalf("Failed to listen on port %d: %v", port, err)
	}
	defer l.Close()

	// Now it is in use, so IsPortOpen should be true
	if !IsPortOpen(port) {
		t.Errorf("Port %d should be in use (IsPortOpen=true), but got false", port)
	}
}

func TestScanRange(t *testing.T) {
	// Find two ports
	port1 := getFreePort(t)
	// We need another one. getFreePort binds temporarily then releases.
	// But we need to hold them for the test.

	// Let's bind port1
	l1, err := net.Listen("tcp", fmt.Sprintf(":%d", port1))
	if err != nil {
		t.Fatalf("Failed to bind port1 %d: %v", port1, err)
	}
	defer l1.Close()

	// Bind port2 (assuming next one is free or find another)
	// Simple approach: try binding port1 + 1.
	port2 := port1 + 1
	l2, err := net.Listen("tcp", fmt.Sprintf(":%d", port2))
	if err != nil {
		// If fails, we skip testing exactly two ports or find another strategy
		t.Logf("Could not bind port2 %d, skipping part of test: %v", port2, err)
		port2 = 0
	} else {
		defer l2.Close()
	}

	// Scan range covering these ports
	// We scan [port1, port1+1]
	start := port1
	end := port1
	if port2 > 0 {
		end = port2
	}

	openPorts := ScanRange(start, end)

	if len(openPorts) == 0 {
		t.Error("Expected at least one open port")
	}

	found1 := false
	found2 := false
	for _, p := range openPorts {
		if p == port1 {
			found1 = true
		}
		if port2 > 0 && p == port2 {
			found2 = true
		}
	}

	if !found1 {
		t.Errorf("Did not find port1 %d in scan results", port1)
	}
	if port2 > 0 && !found2 {
		t.Errorf("Did not find port2 %d in scan results", port2)
	}
}

func TestScanCommon(t *testing.T) {
	// We can't guarantee common ports (3000, 8080 etc) are open or closed.
	// But we can test it runs without panic.
	ports := ScanCommon()
	t.Logf("Common open ports found: %v", ports)
}
