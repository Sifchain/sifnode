package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func listenForSignals(cmd *exec.Cmd) {
	// Set up channel to listen for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-sigChan

	// Stop the process when a signal is received
	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			log.Fatalf("Failed to kill process: %v", err)
		}
		log.Println("Process killed successfully")
	}
}
