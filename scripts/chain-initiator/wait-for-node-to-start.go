package main

import (
	"log"
	"time"
)

func waitForNodeToStart(node string) {
	timeout := 60 * time.Second
	start := time.Now()

	// Wait for the node to be running with timout
	for !isNodeRunning(node) {
		if time.Since(start) > timeout {
			log.Fatalf("Node did not start within the specified timeout")
		}
		log.Println("Waiting for node to start...")
		time.Sleep(5 * time.Second)
	}
	log.Println("Node is running.")
}
