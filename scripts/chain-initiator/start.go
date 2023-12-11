package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func start(cmdPath, homePath string) *exec.Cmd {
	// Command and arguments
	args := []string{"start", "--home", homePath}

	// Set up the command
	cmd := exec.Command(cmdPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and stream the output in a goroutine to avoid blocking
	go func() {
		err := cmd.Run()
		if err != nil {
			// Check if the error is because of the process being killed
			if exitErr, ok := err.(*exec.ExitError); ok {
				// If the process was killed, log it as a non-fatal error
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
					log.Printf("Process was killed: %v", err)
					return
				}
			}
			// For other errors, log them as fatal
			log.Fatalf("Command execution failed: %v", err)
		}
	}()

	return cmd
}
