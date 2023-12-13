package main

import (
	"net"
	"net/http"
	"strings"
)

func isNodeRunning(node string) bool {
	// Remove the "tcp://" prefix if present
	if strings.HasPrefix(node, "tcp://") {
		node = strings.TrimPrefix(node, "tcp://")
	}

	// Attempt to make a TCP connection
	conn, err := net.Dial("tcp", node)
	if err == nil {
		conn.Close()
		return true
	}

	// If TCP connection fails, attempt an HTTP GET request
	resp, err := http.Get("http://" + node)
	if err == nil {
		resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}

	return false
}
