package ports

import (
	"fmt"
	"net"
)

// IsPortOpen checks if a port is currently in use
func IsPortOpen(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return true // Port in use
	}
	listener.Close()
	return false
}

// ScanRange scans a range of ports and returns those that are open
func ScanRange(start, end int) []int {
	var openPorts []int
	for port := start; port <= end; port++ {
		if IsPortOpen(port) {
			openPorts = append(openPorts, port)
		}
	}
	return openPorts
}

// ScanCommon scans commonly used development ports
func ScanCommon() []int {
	commonPorts := []int{3000, 3001, 4200, 5173, 8080, 8000, 8888, 9000}
	var openPorts []int
	for _, port := range commonPorts {
		if IsPortOpen(port) {
			openPorts = append(openPorts, port)
		}
	}
	return openPorts
}
